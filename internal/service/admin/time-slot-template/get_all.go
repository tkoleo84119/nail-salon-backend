package adminTimeSlotTemplate

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, req adminTimeSlotTemplateModel.GetAllParsedRequest) (*adminTimeSlotTemplateModel.GetAllResponse, error) {
	total, timeSlotTemplates, err := s.repo.Template.GetAllTimeSlotTemplatesByFilter(ctx, sqlxRepo.GetAllTimeSlotTemplatesByFilterParams{
		Name:   req.Name,
		Limit:  &req.Limit,
		Offset: &req.Offset,
		Sort:   &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot template list", err)
	}

	response := &adminTimeSlotTemplateModel.GetAllResponse{
		Total: total,
		Items: make([]adminTimeSlotTemplateModel.GetAllItem, len(timeSlotTemplates)),
	}

	for i, item := range timeSlotTemplates {
		response.Items[i] = adminTimeSlotTemplateModel.GetAllItem{
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
