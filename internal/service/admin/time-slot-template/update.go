package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotTemplateService struct {
	queries dbgen.Querier
	repo    *sqlx.Repositories
}

func NewUpdateTimeSlotTemplateService(queries dbgen.Querier, repo *sqlx.Repositories) *UpdateTimeSlotTemplateService {
	return &UpdateTimeSlotTemplateService{
		queries: queries,
		repo:    repo,
	}
}

func (s *UpdateTimeSlotTemplateService) UpdateTimeSlotTemplate(ctx context.Context, templateID string, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error) {
	// Validate at least one field is provided
	if !req.HasUpdate() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Parse template ID
	templateIDInt, err := utils.ParseID(templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid template ID", err)
	}

	// Check if template exists
	_, err = s.queries.GetTimeSlotTemplateByID(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotTemplateNotFound, "time slot template not found", err)
	}

	// Update the template using sqlx repository
	response, err := s.repo.Template.UpdateTimeSlotTemplate(ctx, templateIDInt, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update time slot template", err)
	}

	return response, nil
}
