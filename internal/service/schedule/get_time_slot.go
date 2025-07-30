package schedule

import (
	"context"
	"database/sql"
	"errors"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetTimeSlotService struct {
	queries *dbgen.Queries
}

func NewGetTimeSlotService(queries *dbgen.Queries) *GetTimeSlotService {
	return &GetTimeSlotService{
		queries: queries,
	}
}

func (s *GetTimeSlotService) GetTimeSlotsBySchedule(ctx context.Context, scheduleID string, customerContext common.CustomerContext) (*scheduleModel.GetTimeSlotsByScheduleResponse, error) {
	// Input validation & ID parsing
	scheduleIDInt, err := utils.ParseID(scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid schedule ID", err)
	}

	// Check if customer is blacklisted
	customer, err := s.queries.GetCustomerByID(ctx, customerContext.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer", err)
	}

	// If customer is blacklisted, return empty array
	if customer.IsBlacklisted.Bool {
		return &scheduleModel.GetTimeSlotsByScheduleResponse{
			Items: []scheduleModel.TimeSlotResponseItem{},
		}, nil
	}

	// Verify schedule exists
	_, err = s.queries.GetScheduleByID(ctx, scheduleIDInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &scheduleModel.GetTimeSlotsByScheduleResponse{
				Items: []scheduleModel.TimeSlotResponseItem{},
			}, nil
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get schedule", err)
	}

	// Get available time slots for the schedule
	timeSlots, err := s.queries.GetAvailableTimeSlotsByScheduleID(ctx, scheduleIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get time slots", err)
	}

	// Build response
	items := make([]scheduleModel.TimeSlotResponseItem, 0, len(timeSlots))
	for _, slot := range timeSlots {
		startTime := utils.PgTimeToTime(slot.StartTime)
		endTime := utils.PgTimeToTime(slot.EndTime)
		durationMinutes := int(endTime.Sub(startTime).Minutes())

		items = append(items, scheduleModel.TimeSlotResponseItem{
			ID:              utils.FormatID(slot.ID),
			StartTime:       utils.PgTimeToTimeString(slot.StartTime),
			EndTime:         utils.PgTimeToTimeString(slot.EndTime),
			DurationMinutes: durationMinutes,
		})
	}

	return &scheduleModel.GetTimeSlotsByScheduleResponse{
		Items: items,
	}, nil
}
