package adminTimeSlotTemplate

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
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

func (s *CreateTimeSlotTemplateItemService) CreateTimeSlotTemplateItem(ctx context.Context, templateID string, req adminTimeSlotTemplateModel.CreateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.CreateTimeSlotTemplateItemResponse, error) {
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
		StartTime:  utils.TimeToPgTime(startTime),
		EndTime:    utils.TimeToPgTime(endTime),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create time slot template item", err)
	}

	return &adminTimeSlotTemplateModel.CreateTimeSlotTemplateItemResponse{
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
		existingStart := utils.PgTimeToTime(item.StartTime)
		existingEnd := utils.PgTimeToTime(item.EndTime)

		// Check if new slot overlaps with existing slot
		if startTime.Before(existingEnd) && endTime.After(existingStart) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotConflict)
		}
	}
	return nil
}
