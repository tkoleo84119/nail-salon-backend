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

type UpdateTimeSlotTemplateItemServiceInterface interface {
	UpdateTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, req timeSlotTemplate.UpdateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateItemResponse, error)
}

type UpdateTimeSlotTemplateItemService struct {
	queries dbgen.Querier
}

func NewUpdateTimeSlotTemplateItemService(queries dbgen.Querier) *UpdateTimeSlotTemplateItemService {
	return &UpdateTimeSlotTemplateItemService{
		queries: queries,
	}
}

func (s *UpdateTimeSlotTemplateItemService) UpdateTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, req timeSlotTemplate.UpdateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateItemResponse, error) {
	// Parse template ID
	templateIDInt, err := utils.ParseID(templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid template ID", err)
	}

	// Parse item ID
	itemIDInt, err := utils.ParseID(itemID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid item ID", err)
	}

	// Check if template exists
	_, err = s.queries.GetTimeSlotTemplateByID(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotTemplateNotFound, "time slot template not found", err)
	}

	// Check if item exists and belongs to the template
	existingItem, err := s.queries.GetTimeSlotTemplateItemByID(ctx, itemIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotTemplateItemNotFound, "time slot template item not found", err)
	}

	// Verify the item belongs to the specified template
	if existingItem.TemplateID != templateIDInt {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotTemplateItemNotFound, "time slot template item not found in specified template", nil)
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

	// Get other existing template items to check for conflicts (excluding current item)
	otherItems, err := s.queries.GetTimeSlotTemplateItemsByTemplateIDExcluding(ctx, dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams{
		TemplateID: templateIDInt,
		ID:         itemIDInt,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get existing template items", err)
	}

	// Check for time conflicts with other existing items
	if err := s.checkTimeConflicts(startTime, endTime, otherItems); err != nil {
		return nil, err
	}

	// Update time slot template item
	updatedItem, err := s.queries.UpdateTimeSlotTemplateItem(ctx, dbgen.UpdateTimeSlotTemplateItemParams{
		ID:        itemIDInt,
		StartTime: pgtype.Time{Microseconds: int64(startTime.Hour()*3600+startTime.Minute()*60+startTime.Second()) * 1000000, Valid: true},
		EndTime:   pgtype.Time{Microseconds: int64(endTime.Hour()*3600+endTime.Minute()*60+endTime.Second()) * 1000000, Valid: true},
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot template item", err)
	}

	return &timeSlotTemplate.UpdateTimeSlotTemplateItemResponse{
		ID:         utils.FormatID(updatedItem.ID),
		TemplateID: utils.FormatID(updatedItem.TemplateID),
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}, nil
}

// checkTimeConflicts validates that the updated time slot doesn't overlap with other existing items
func (s *UpdateTimeSlotTemplateItemService) checkTimeConflicts(startTime, endTime time.Time, existingItems []dbgen.TimeSlotTemplateItem) error {
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