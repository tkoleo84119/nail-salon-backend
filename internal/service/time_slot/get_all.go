package timeSlot

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	timeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
}

func NewGetAll(queries *dbgen.Queries) *GetAll {
	return &GetAll{
		queries: queries,
	}
}

func (s *GetAll) GetAll(ctx context.Context, scheduleID int64, isBlacklisted bool) (*timeSlotModel.GetAllResponse, error) {
	timeSlots := make([]timeSlotModel.GetAllResponseItem, 0)
	response := timeSlotModel.GetAllResponse{
		TimeSlots: timeSlots,
	}

	// If customer is blacklisted, return empty array
	if isBlacklisted {
		return &response, nil
	}

	// Verify schedule exists
	exists, err := s.queries.CheckScheduleExistsByID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get schedule", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	}

	// Get available time slots for the schedule
	rawTimeSlots, err := s.queries.GetAvailableTimeSlotsByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get time slots", err)
	}

	for _, rawTimeSlot := range rawTimeSlots {
		timeSlots = append(timeSlots, timeSlotModel.GetAllResponseItem{
			ID:              utils.FormatID(rawTimeSlot.ID),
			StartTime:       utils.PgTimeToTimeString(rawTimeSlot.StartTime),
			EndTime:         utils.PgTimeToTimeString(rawTimeSlot.EndTime),
			DurationMinutes: int(utils.PgTimeToTime(rawTimeSlot.EndTime).Sub(utils.PgTimeToTime(rawTimeSlot.StartTime)).Minutes()),
		})
	}

	response.TimeSlots = timeSlots

	return &response, nil
}
