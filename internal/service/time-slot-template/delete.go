package timeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteTimeSlotTemplateService struct {
	queries dbgen.Querier
}

func NewDeleteTimeSlotTemplateService(queries dbgen.Querier) *DeleteTimeSlotTemplateService {
	return &DeleteTimeSlotTemplateService{
		queries: queries,
	}
}

func (s *DeleteTimeSlotTemplateService) DeleteTimeSlotTemplate(ctx context.Context, templateID string, staffContext common.StaffContext) (*timeSlotTemplate.DeleteTimeSlotTemplateResponse, error) {
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

	// Delete the time slot template (cascade will delete items)
	err = s.queries.DeleteTimeSlotTemplate(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete time slot template", err)
	}

	return &timeSlotTemplate.DeleteTimeSlotTemplateResponse{
		Deleted: []string{templateID},
	}, nil
}