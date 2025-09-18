package adminTimeSlot

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, scheduleID int64, timeSlotID int64, req adminTimeSlotModel.UpdateRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminTimeSlotModel.UpdateResponse, error) {
	// Validate that both start time and end time are provided together
	if (req.StartTime != nil && req.EndTime == nil) || (req.StartTime == nil && req.EndTime != nil) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotCannotUpdateSeparately)
	}

	// Get time slot information
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
	}

	// Verify time slot belongs to the specified schedule
	if timeSlot.ScheduleID != scheduleID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotBelongToSchedule)
	}

	// Get schedule information
	scheduleInfo, err := s.queries.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	}

	// check if schedule is past date
	if err := s.checkScheduleIsNotPastDate(scheduleInfo.WorkDate); err != nil {
		return nil, err
	}

	// Get stylist information
	stylist, err := s.queries.GetStylistByID(ctx, scheduleInfo.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only modify their own schedules
	if updaterRole == common.RoleStylist {
		if stylist.StaffUserID != updaterID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if staff has access to this store
	if err := utils.CheckStoreAccess(scheduleInfo.StoreID, updaterStoreIDs, updaterRole); err != nil {
		return nil, err
	}

	if req.StartTime != nil && req.EndTime != nil {
		startTimeParsed, err := utils.TimeStringToTime(*req.StartTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid start time format", err)
		}
		endTimeParsed, err := utils.TimeStringToTime(*req.EndTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid end time format", err)
		}

		// Validate time range
		if !endTimeParsed.After(startTimeParsed) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
		}

		hasOverlap, err := s.queries.CheckTimeSlotOverlapExcluding(ctx, dbgen.CheckTimeSlotOverlapExcludingParams{
			ScheduleID: scheduleID,
			ID:         timeSlotID,
			Column3:    utils.TimePtrToPgTime(&startTimeParsed),
			Column4:    utils.TimePtrToPgTime(&endTimeParsed),
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot overlap", err)
		}
		if hasOverlap {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}

	// Update time slot using sqlx repository
	response, err := s.repo.TimeSlot.UpdateTimeSlot(ctx, timeSlotID, sqlxRepo.UpdateTimeSlotParams{
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot", err)
	}

	return &adminTimeSlotModel.UpdateResponse{
		ID:          utils.FormatID(response.ID),
		ScheduleID:  utils.FormatID(response.ScheduleID),
		StartTime:   utils.PgTimeToTimeString(response.StartTime),
		EndTime:     utils.PgTimeToTimeString(response.EndTime),
		IsAvailable: utils.PgBoolToBool(response.IsAvailable),
	}, nil
}

func (s *Update) checkScheduleIsNotPastDate(workDatePg pgtype.Date) error {
	workDate := workDatePg.Time
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to load location", err)
	}

	now := time.Now().In(loc)
	if workDate.Before(now) {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotUpdateBeforeToday)
	}

	return nil
}
