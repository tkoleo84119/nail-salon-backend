package timeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteTimeSlotTemplateItemServiceInterface interface {
	DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*timeSlotTemplate.DeleteTimeSlotTemplateItemResponse, error)
}

type DeleteTimeSlotTemplateItemService struct {
	queries dbgen.Querier
}

func NewDeleteTimeSlotTemplateItemService(queries dbgen.Querier) *DeleteTimeSlotTemplateItemService {
	return &DeleteTimeSlotTemplateItemService{
		queries: queries,
	}
}

func (s *DeleteTimeSlotTemplateItemService) DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*timeSlotTemplate.DeleteTimeSlotTemplateItemResponse, error) {
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

	// Delete the time slot template item
	err = s.queries.DeleteTimeSlotTemplateItem(ctx, itemIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot template item", err)
	}

	return &timeSlotTemplate.DeleteTimeSlotTemplateItemResponse{
		Deleted: []string{itemID},
	}, nil
}