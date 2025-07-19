package timeSlotTemplate

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateTimeSlotTemplateItemService struct {
	queries dbgen.Querier
}

func NewCreateTimeSlotTemplateItemService(queries dbgen.Querier) *CreateTimeSlotTemplateItemService {
	return &CreateTimeSlotTemplateItemService{
		queries: queries,
	}
}

func (s *CreateTimeSlotTemplateItemService) CreateTimeSlotTemplateItem(ctx context.Context, templateID string, req timeSlotTemplate.CreateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateItemResponse, error) {
	// Parse template ID
	templateIDInt, err := utils.ParseID(templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid template ID", err)
	}

	// Check if template exists
	_, err = s.queries.GetTimeSlotTemplateByID(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotTemplateNotFound, "time slot template not found", err)
	}

	// Validate time format and range
	startTime, err := common.ParseTimeSlot(req.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", err)
	}

	endTime, err := common.ParseTimeSlot(req.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end time format", err)
	}

	// Validate time range
	if !endTime.After(startTime) {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "end time must be after start time", nil)
	}

	// Get existing template items to check for conflicts
	existingItems, err := s.queries.GetTimeSlotTemplateItemsByTemplateID(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get existing template items", err)
	}

	// Check for time conflicts with existing items
	if err := s.checkTimeConflicts(startTime, endTime, existingItems); err != nil {
		return nil, err
	}

	// Generate item ID
	itemID := utils.GenerateID()

	// Create time slot template item
	item, err := s.queries.CreateTimeSlotTemplateItem(ctx, dbgen.CreateTimeSlotTemplateItemParams{
		ID:         itemID,
		TemplateID: templateIDInt,
		StartTime:  pgtype.Time{Microseconds: int64(startTime.Hour()*3600+startTime.Minute()*60+startTime.Second()) * 1000000, Valid: true},
		EndTime:    pgtype.Time{Microseconds: int64(endTime.Hour()*3600+endTime.Minute()*60+endTime.Second()) * 1000000, Valid: true},
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot template item", err)
	}

	return &timeSlotTemplate.CreateTimeSlotTemplateItemResponse{
		ID:         utils.FormatID(item.ID),
		TemplateID: utils.FormatID(item.TemplateID),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}, nil
}

// checkTimeConflicts validates that the new time slot doesn't overlap with existing items
func (s *CreateTimeSlotTemplateItemService) checkTimeConflicts(startTime, endTime time.Time, existingItems []dbgen.TimeSlotTemplateItem) error {
	for _, item := range existingItems {
		// Convert pgtype.Time to time.Time for comparison
		existingStart := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(item.StartTime.Microseconds) * time.Microsecond)
		existingEnd := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(item.EndTime.Microseconds) * time.Microsecond)

		// Check if new slot overlaps with existing slot
		if startTime.Before(existingEnd) && endTime.After(existingStart) {
			return errorCodes.NewServiceError(errorCodes.ScheduleTimeConflict, "time slot overlaps with existing template item", nil)
		}
	}
	return nil
}
