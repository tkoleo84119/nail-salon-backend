package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetTimeSlotTemplateListServiceInterface defines the interface for getting time slot template list
type GetTimeSlotTemplateListServiceInterface interface {
	GetTimeSlotTemplateList(ctx context.Context, req adminTimeSlotTemplateModel.GetTimeSlotTemplateListRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse, error)
}

type GetTimeSlotTemplateListService struct {
	queries              *dbgen.Queries
	timeSlotTemplateRepo sqlxRepo.TimeSlotTemplateRepositoryInterface
}

func NewGetTimeSlotTemplateListService(queries *dbgen.Queries, timeSlotTemplateRepo sqlxRepo.TimeSlotTemplateRepositoryInterface) *GetTimeSlotTemplateListService {
	return &GetTimeSlotTemplateListService{
		queries:              queries,
		timeSlotTemplateRepo: timeSlotTemplateRepo,
	}
}

func (s *GetTimeSlotTemplateListService) GetTimeSlotTemplateList(ctx context.Context, req adminTimeSlotTemplateModel.GetTimeSlotTemplateListRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse, error) {
	timeSlotTemplates, total, err := s.timeSlotTemplateRepo.GetTimeSlotTemplateList(ctx, sqlxRepo.GetTimeSlotTemplateListParams{
		Name:   req.Name,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot template list", err)
	}

	response := &adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse{
		Total: total,
		Items: make([]adminTimeSlotTemplateModel.GetTimeSlotTemplateListItem, len(timeSlotTemplates)),
	}
	for i, item := range timeSlotTemplates {
		response.Items[i] = adminTimeSlotTemplateModel.GetTimeSlotTemplateListItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			Note:      item.Note,
			Updater:   utils.FormatID(item.Updater),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return response, nil
}
