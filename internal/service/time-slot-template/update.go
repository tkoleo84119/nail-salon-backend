package timeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotTemplateService struct {
	queries              dbgen.Querier
	timeSlotTemplateRepo sqlx.TimeSlotTemplateRepositoryInterface
}

func NewUpdateTimeSlotTemplateService(queries dbgen.Querier, timeSlotTemplateRepo sqlx.TimeSlotTemplateRepositoryInterface) *UpdateTimeSlotTemplateService {
	return &UpdateTimeSlotTemplateService{
		queries:              queries,
		timeSlotTemplateRepo: timeSlotTemplateRepo,
	}
}

func (s *UpdateTimeSlotTemplateService) UpdateTimeSlotTemplate(ctx context.Context, templateID string, req timeSlotTemplate.UpdateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateResponse, error) {
	// Validate at least one field is provided
	if !req.HasUpdate() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

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

	// Update the template using sqlx repository
	response, err := s.timeSlotTemplateRepo.UpdateTimeSlotTemplate(ctx, templateIDInt, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot template", err)
	}

	return response, nil
}
