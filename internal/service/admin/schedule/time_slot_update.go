package adminSchedule

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotService struct {
	queries dbgen.Querier
	repo    *sqlx.Repositories
}

func NewUpdateTimeSlotService(queries dbgen.Querier, repo *sqlx.Repositories) *UpdateTimeSlotService {
	return &UpdateTimeSlotService{
		queries: queries,
		repo:    repo,
	}
}

func (s *UpdateTimeSlotService) UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req adminScheduleModel.UpdateTimeSlotRequest, staffContext common.StaffContext) (*adminScheduleModel.UpdateTimeSlotResponse, error) {
	// Validate at least one field is provided
	if !req.HasUpdate() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Validate that both start time and end time are provided together
	if (req.StartTime != nil && req.EndTime == nil) || (req.StartTime == nil && req.EndTime != nil) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotCannotUpdateSeparately)
	}

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

	// Check if time slot is booked (cannot update booked time slots)
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotAlreadyBookedDoNotUpdate)
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

	if req.StartTime != nil && req.EndTime != nil {
		startTimeParsed, err := utils.TimeStringToTime(*req.StartTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid start time format", err)
		}
		endTimeParsed, err := utils.TimeStringToTime(*req.EndTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid end time format", err)
		}

		// Validate time range
		if !endTimeParsed.After(startTimeParsed) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotInvalidTimeRange)
		}

		hasOverlap, err := s.queries.CheckTimeSlotOverlapExcluding(ctx, dbgen.CheckTimeSlotOverlapExcludingParams{
			ScheduleID: scheduleIDInt,
			ID:         timeSlotIDInt,
			StartTime:  utils.TimeToPgTime(startTimeParsed),
			EndTime:    utils.TimeToPgTime(endTimeParsed),
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot overlap", err)
		}
		if hasOverlap {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}

	// Update time slot using sqlx repository
	response, err := s.repo.TimeSlot.UpdateTimeSlot(ctx, timeSlotIDInt, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot", err)
	}

	// Set the correct schedule ID in response
	response.ScheduleID = scheduleID

	return response, nil
}
