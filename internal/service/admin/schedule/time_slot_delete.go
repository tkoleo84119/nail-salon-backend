package adminSchedule

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteTimeSlotService struct {
	queries dbgen.Querier
}

func NewDeleteTimeSlotService(queries dbgen.Querier) *DeleteTimeSlotService {
	return &DeleteTimeSlotService{
		queries: queries,
	}
}

func (s *DeleteTimeSlotService) DeleteTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, staffContext common.StaffContext) (*adminScheduleModel.DeleteTimeSlotResponse, error) {
	// Parse IDs
	scheduleIDInt, err := utils.ParseID(scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid schedule ID", err)
	}

	timeSlotIDInt, err := utils.ParseID(timeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid time slot ID", err)
	}

	// Get time slot information
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot not found", err)
	}

	// Verify time slot belongs to the specified schedule
	if timeSlot.ScheduleID != scheduleIDInt {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotBelongToSchedule)
	}

	// Check if time slot is booked (cannot delete booked time slots)
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotAlreadyBookedDoNotDelete)
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

	// Check permission: STYLIST can only delete their own schedules
	if staffContext.Role == adminStaffModel.RoleStylist {
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
		return nil, errorCodes.NewServiceError(errorCodes.StaffStoreNotFound, "store not found", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.StaffStoreNotActive, "store is not active", err)
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

	// Delete time slot
	err = s.queries.DeleteTimeSlotByID(ctx, timeSlotIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot", err)
	}

	return &adminScheduleModel.DeleteTimeSlotResponse{
		Deleted: timeSlotID,
	}, nil
}
