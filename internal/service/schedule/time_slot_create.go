package schedule

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateTimeSlotService struct {
	queries dbgen.Querier
}

func NewCreateTimeSlotService(queries dbgen.Querier) *CreateTimeSlotService {
	return &CreateTimeSlotService{
		queries: queries,
	}
}

func (s *CreateTimeSlotService) CreateTimeSlot(ctx context.Context, scheduleID string, req schedule.CreateTimeSlotRequest, staffContext common.StaffContext) (*schedule.CreateTimeSlotResponse, error) {
	// Parse schedule ID
	scheduleIDInt, err := utils.ParseID(scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid schedule ID", err)
	}

	// Validate time format
	startTime, err := utils.TimeStringToTime(req.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
	}
	endTime, err := utils.TimeStringToTime(req.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
	}

	// Validate time range
	if !endTime.After(startTime) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotInvalidTimeRange)
	}

	// Get schedule information
	scheduleInfo, err := s.queries.GetScheduleByID(ctx, scheduleIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotFound, "schedule not found", err)
	}

	// Get stylist information
	stylist, err := s.queries.GetStylistByID(ctx, scheduleInfo.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StylistNotFound, "stylist not found", err)
	}

	// Check permission: STYLIST can only modify their own schedules
	if staffContext.Role == staff.RoleStylist {
		staffUserID, err := utils.ParseID(staffContext.UserID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
		}
		if stylist.StaffUserID != staffUserID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if store exists and is active
	store, err := s.queries.GetStoreByID(ctx, scheduleInfo.StoreID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.UserStoreNotFound, "store not found", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.UserStoreNotActive, "store is not active", err)
	}

	// Check if staff has access to this store
	hasAccess := false
	storeIDStr := utils.FormatID(scheduleInfo.StoreID)
	for _, storeAccess := range staffContext.StoreList {
		if storeAccess.ID == storeIDStr {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check for time slot overlap
	hasOverlap, err := s.queries.CheckTimeSlotOverlap(ctx, dbgen.CheckTimeSlotOverlapParams{
		ScheduleID: scheduleIDInt,
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
		ScheduleID: scheduleIDInt,
		StartTime:  utils.TimeToPgTime(startTime),
		EndTime:    utils.TimeToPgTime(endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot", err)
	}

	// Build response
	response := &schedule.CreateTimeSlotResponse{
		ID:          utils.FormatID(createdTimeSlot.ID),
		ScheduleID:  scheduleID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: createdTimeSlot.IsAvailable.Bool,
	}

	return response, nil
}
