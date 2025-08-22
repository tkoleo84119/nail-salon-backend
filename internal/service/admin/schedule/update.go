package adminSchedule

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID int64, scheduleID int64, req adminScheduleModel.UpdateParsedRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminScheduleModel.UpdateResponse, error) {
	// Check store access for the updater (except SUPER_ADMIN)
	if updaterRole != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, updaterStoreIDs); err != nil {
			return nil, err
		}
	}

	// Verify stylist exists
	exist, err := s.queries.CheckStylistExistAndActive(ctx, req.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check stylist existence", err)
	}
	if !exist {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Get current schedule to verify ownership and existence
	currentSchedule, err := s.queries.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule", err)
	}

	// Verify schedule belongs to the specified store
	if currentSchedule.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStore)
	}
	if currentSchedule.StylistID != req.StylistID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotBelongToStylist)
	}

	// Check if user can update this schedule
	if updaterRole == common.RoleStylist {
		if currentSchedule.StylistID != updaterID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// if update date, check if schedule can be updated
	var workDate time.Time
	if req.WorkDate != nil {
		canUpdate, err := s.queries.CheckScheduleCanUpdateDate(ctx, scheduleID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check schedule can update date", err)
		}
		if !canUpdate {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyBookedDoNotUpdateDate)
		}

		workDate = *req.WorkDate
		// Check if the new date is not in the past
		today := time.Now().Truncate(24 * time.Hour)
		if workDate.Before(today) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotCreateBeforeToday)
		}

		// Check if schedule with new date already exists (if date is being changed)
		if workDate != currentSchedule.WorkDate.Time {
			exists, err := s.queries.CheckScheduleDateExists(ctx, dbgen.CheckScheduleDateExistsParams{
				StoreID:   storeID,
				StylistID: req.StylistID,
				WorkDate:  utils.TimeToPgDate(workDate),
			})
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check schedule existence", err)
			}
			if exists {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyExists)
			}
		}
	}

	// Update schedule
	updatedSchedule, err := s.repo.Schedule.UpdateSchedule(ctx, scheduleID, sqlxRepo.UpdateScheduleParams{
		WorkDate: utils.TimeToPgDate(workDate),
		Note:     req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to update schedule", err)
	}

	// Build response
	response := &adminScheduleModel.UpdateResponse{
		ID: utils.FormatID(updatedSchedule.ID),
	}

	return response, nil
}
