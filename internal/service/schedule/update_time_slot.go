package schedule

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotServiceInterface interface {
	UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req schedule.UpdateTimeSlotRequest, staffContext common.StaffContext) (*schedule.UpdateTimeSlotResponse, error)
}

type UpdateTimeSlotService struct {
	queries      dbgen.Querier
	timeSlotRepo sqlx.TimeSlotRepositoryInterface
}

func NewUpdateTimeSlotService(queries dbgen.Querier, timeSlotRepo sqlx.TimeSlotRepositoryInterface) *UpdateTimeSlotService {
	return &UpdateTimeSlotService{
		queries:      queries,
		timeSlotRepo: timeSlotRepo,
	}
}

func (s *UpdateTimeSlotService) UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req schedule.UpdateTimeSlotRequest, staffContext common.StaffContext) (*schedule.UpdateTimeSlotResponse, error) {
	// Validate at least one field is provided
	if req.StartTime == nil && req.EndTime == nil && req.IsAvailable == nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Validate that both start time and end time are provided together
	if (req.StartTime != nil && req.EndTime == nil) || (req.StartTime == nil && req.EndTime != nil) {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotCannotUpdateSeparately, "start time and end time must be provided together", nil)
	}

	// Parse IDs
	scheduleIDInt, err := utils.ParseID(scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid schedule ID", err)
	}

	timeSlotIDInt, err := utils.ParseID(timeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid time slot ID", err)
	}

	// Get time slot information
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot not found", err)
	}

	// Verify time slot belongs to the specified schedule
	if timeSlot.ScheduleID != scheduleIDInt {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot does not belong to specified schedule", nil)
	}

	// Check if time slot is booked (cannot update booked time slots)
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotAlreadyBookedDoNotUpdate, "cannot update booked time slot", nil)
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
		if stylist.StaffUserID.Int64 != staffUserID {
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

	var startTime, endTime *pgtype.Time
	if req.StartTime != nil && req.EndTime != nil {
		startTimeParsed, err := schedule.ParseTimeSlot(*req.StartTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
		}
		endTimeParsed, err := schedule.ParseTimeSlot(*req.EndTime)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
		}

		// Validate time range
		if !endTimeParsed.After(startTimeParsed) {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "end time must be after start time", nil)
		}

		startTime = &pgtype.Time{Microseconds: int64(startTimeParsed.Hour()*3600+startTimeParsed.Minute()*60+startTimeParsed.Second()) * 1000000, Valid: true}
		endTime = &pgtype.Time{Microseconds: int64(endTimeParsed.Hour()*3600+endTimeParsed.Minute()*60+endTimeParsed.Second()) * 1000000, Valid: true}

		hasOverlap, err := s.queries.CheckTimeSlotOverlapExcluding(ctx, dbgen.CheckTimeSlotOverlapExcludingParams{
			ScheduleID: scheduleIDInt,
			ID:         timeSlotIDInt,
			StartTime:  *startTime,
			EndTime:    *endTime,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check time slot overlap", err)
		}
		if hasOverlap {
			return nil, errorCodes.NewServiceError(errorCodes.ScheduleTimeConflict, "time slot overlaps with existing time slots", nil)
		}
	}

	// Update time slot using sqlx repository
	response, err := s.timeSlotRepo.UpdateTimeSlot(ctx, timeSlotIDInt, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot", err)
	}

	// Set the correct schedule ID in response
	response.ScheduleID = scheduleID

	return response, nil
}
