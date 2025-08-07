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
	// Check at least one field is provided
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Verify stylist exists
	_, err := s.queries.GetStylistByID(ctx, req.StylistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist", err)
	}

	// Check store access for the updater (except SUPER_ADMIN)
	if updaterRole != common.RoleSuperAdmin {
		hasAccess, err := utils.CheckStoreAccess(storeID, updaterStoreIDs)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
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

	// Get updated schedule with time slots
	rows, err := s.queries.GetScheduleWithTimeSlotsByID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get updated schedule with time slots", err)
	}

	// check all time slots are available
	for _, row := range rows {
		if !row.IsAvailable.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleAlreadyBookedDoNotUpdate)
		}
	}

	// Authorization: Check if user can update this schedule
	if updaterRole == common.RoleStylist {
		if currentSchedule.StylistID != updaterID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Prepare update parameters
	if req.WorkDate != nil {
		// Parse and validate date format
		parsedDate, err := time.Parse("2006-01-02", *req.WorkDate)
		if err != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValFieldDateFormat)
		}

		// Check if the new date is not in the past
		today := time.Now().Truncate(24 * time.Hour)
		if parsedDate.Before(today) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotCreateBeforeToday)
		}

		// Check if schedule with new date already exists (if date is being changed)
		if parsedDate != currentSchedule.WorkDate.Time {
			exists, err := s.queries.CheckScheduleExists(ctx, dbgen.CheckScheduleExistsParams{
				StoreID:   storeID,
				StylistID: req.StylistID,
				WorkDate:  utils.TimeToPgDate(parsedDate),
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
		WorkDate: req.WorkDate,
		Note:     req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to update schedule", err)
	}

	// Build response
	response := &adminScheduleModel.UpdateResponse{
		ID:        utils.FormatID(updatedSchedule.ID),
		WorkDate:  utils.PgDateToDateString(updatedSchedule.WorkDate),
		Note:      utils.PgTextToString(updatedSchedule.Note),
		TimeSlots: make([]adminScheduleModel.UpdateTimeSlotInfo, len(rows)),
	}

	for i, row := range rows {
		if row.TimeSlotID.Valid {
			timeSlot := adminScheduleModel.UpdateTimeSlotInfo{
				ID:          utils.FormatID(row.TimeSlotID.Int64),
				StartTime:   utils.PgTimeToTimeString(row.StartTime),
				EndTime:     utils.PgTimeToTimeString(row.EndTime),
				IsAvailable: utils.PgBoolToBool(row.IsAvailable),
			}
			response.TimeSlots[i] = timeSlot
		}
	}

	return response, nil
}
