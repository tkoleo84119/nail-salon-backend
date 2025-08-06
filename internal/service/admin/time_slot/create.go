package adminTimeSlot

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries dbgen.Querier
}

func NewCreate(queries dbgen.Querier) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, scheduleID int64, req adminTimeSlotModel.CreateRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminTimeSlotModel.CreateResponse, error) {
	// Validate time format
	startTime, err := utils.TimeStringToTime(req.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid start time format", err)
	}
	endTime, err := utils.TimeStringToTime(req.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValFieldTimeFormat, "invalid end time format", err)
	}

	// Validate time range
	if !endTime.After(startTime) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
	}

	// Get schedule information
	scheduleInfo, err := s.queries.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	}

	// Check if schedule date is not before
	if scheduleInfo.WorkDate.Time.Before(time.Now()) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleCannotCreateBeforeToday)
	}

	// Get stylist information
	stylist, err := s.queries.GetStylistByID(ctx, scheduleInfo.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only modify their own schedules
	if creatorRole == common.RoleStylist {
		if stylist.StaffUserID != creatorID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if store exists and is active
	store, err := s.queries.GetStoreByID(ctx, scheduleInfo.StoreID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	// Check if staff has access to this store
	hasAccess, err := utils.CheckStoreAccess(store.ID, creatorStoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
	}
	if !hasAccess {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check for time slot overlap
	hasOverlap, err := s.queries.CheckTimeSlotOverlap(ctx, dbgen.CheckTimeSlotOverlapParams{
		ScheduleID: scheduleID,
		StartTime:  utils.TimeToPgTime(startTime),
		EndTime:    utils.TimeToPgTime(endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot overlap", err)
	}
	if hasOverlap {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
	}

	// Create time slot
	timeSlotID := utils.GenerateID()
	createdTimeSlot, err := s.queries.CreateTimeSlot(ctx, dbgen.CreateTimeSlotParams{
		ID:         timeSlotID,
		ScheduleID: scheduleID,
		StartTime:  utils.TimeToPgTime(startTime),
		EndTime:    utils.TimeToPgTime(endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot", err)
	}

	// Build response
	response := &adminTimeSlotModel.CreateResponse{
		ID:          utils.FormatID(createdTimeSlot.ID),
		ScheduleID:  utils.FormatID(scheduleID),
		StartTime:   utils.PgTimeToTimeString(createdTimeSlot.StartTime),
		EndTime:     utils.PgTimeToTimeString(createdTimeSlot.EndTime),
		IsAvailable: utils.PgBoolToBool(createdTimeSlot.IsAvailable),
	}

	return response, nil
}
