package adminTimeSlotTemplate

import (
	"context"

	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminTimeSlotTemplateModel.CreateRequest, creatorID int64) (*adminTimeSlotTemplateModel.CreateResponse, error)
}

type DeleteInterface interface {
	Delete(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.DeleteResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminTimeSlotTemplateModel.GetAllParsedRequest) (*adminTimeSlotTemplateModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.GetResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateRequest) (*adminTimeSlotTemplateModel.UpdateResponse, error)
}
