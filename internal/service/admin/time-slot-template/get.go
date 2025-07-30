package adminTimeSlotTemplate

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetTimeSlotTemplateServiceInterface defines the interface for getting a single time slot template
type GetTimeSlotTemplateServiceInterface interface {
	GetTimeSlotTemplate(ctx context.Context, templateID string, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse, error)
}

type GetTimeSlotTemplateService struct {
	queries *dbgen.Queries
}

func NewGetTimeSlotTemplateService(queries *dbgen.Queries) *GetTimeSlotTemplateService {
	return &GetTimeSlotTemplateService{
		queries: queries,
	}
}

func (s *GetTimeSlotTemplateService) GetTimeSlotTemplate(ctx context.Context, templateID string, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse, error) {
	// Parse template ID
	templateIDInt, err := utils.ParseID(templateID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid template ID", err)
	}

	// Get time slot template by ID
	template, err := s.queries.GetTimeSlotTemplateByID(ctx, templateIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot template", err)
	}

	// Get time slot template items
	templateItems, err := s.queries.GetTimeSlotTemplateItemsByTemplateID(ctx, templateIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot template items", err)
	}

	// Build time slot template items response
	items := make([]adminTimeSlotTemplateModel.GetTimeSlotTemplateItemInfo, 0, len(templateItems))
	for _, item := range templateItems {
		items = append(items, adminTimeSlotTemplateModel.GetTimeSlotTemplateItemInfo{
			ID:        utils.FormatID(item.ID),
			StartTime: utils.PgTimeToTimeString(item.StartTime),
			EndTime:   utils.PgTimeToTimeString(item.EndTime),
		})
	}

	// Build response
	response := &adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse{
		ID:        utils.FormatID(template.ID),
		Name:      template.Name,
		Note:      utils.PgTextToString(template.Note),
		Updater:   utils.PgInt8ToIDString(template.Updater),
		CreatedAt: template.CreatedAt.Time.String(),
		UpdatedAt: template.UpdatedAt.Time.String(),
		Items:     items,
	}

	return response, nil
}
