package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetTimeSlotTemplateListService struct {
	repo *sqlxRepo.Repositories
}

func NewGetTimeSlotTemplateListService(repo *sqlxRepo.Repositories) *GetTimeSlotTemplateListService {
	return &GetTimeSlotTemplateListService{
		repo: repo,
	}
}

func (s *GetTimeSlotTemplateListService) GetTimeSlotTemplateList(ctx context.Context, req adminTimeSlotTemplateModel.GetTimeSlotTemplateListParsedRequest) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse, error) {
	total, timeSlotTemplates, err := s.repo.Template.GetAllTimeSlotTemplateByFilter(ctx, sqlxRepo.GetAllTimeSlotTemplateByFilterParams{
		Name:   req.Name,
		Limit:  &req.Limit,
		Offset: &req.Offset,
		Sort:   &req.Sort,
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
			Note:      utils.PgTextToString(item.Note),
			Updater:   utils.PgInt8ToIDString(item.Updater),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	return response, nil
}
