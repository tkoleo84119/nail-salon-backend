package adminTimeSlot

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries dbgen.Querier
}

func NewDelete(queries dbgen.Querier) *Delete {
	return &Delete{
		queries: queries,
	}
}

func (s *Delete) Delete(ctx context.Context, scheduleID int64, timeSlotID int64, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminTimeSlotModel.DeleteResponse, error) {
	// Get time slot information
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot not found", err)
	}
	// Verify time slot belongs to the specified schedule
	if timeSlot.ScheduleID != scheduleID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotBelongToSchedule)
	}
	// Check if time slot is booked (cannot delete booked time slots)
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotAlreadyBookedDoNotDelete)
	}

	// Get schedule information
	scheduleInfo, err := s.queries.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	}

	// Get stylist information
	stylist, err := s.queries.GetStylistByID(ctx, scheduleInfo.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	// Check permission: STYLIST can only delete their own schedules
	if updaterRole == common.RoleStylist {
		if stylist.StaffUserID != updaterID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if staff has access to this store
	if err := utils.CheckStoreAccess(scheduleInfo.StoreID, updaterStoreIDs); err != nil {
		return nil, err
	}

	// Delete time slot
	err = s.queries.DeleteTimeSlotByID(ctx, timeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot", err)
	}

	return &adminTimeSlotModel.DeleteResponse{
		Deleted: utils.FormatID(timeSlotID),
	}, nil
}
