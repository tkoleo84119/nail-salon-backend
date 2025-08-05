package adminTimeSlotTemplate

import (
	"context"

	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
)

type CreateServiceInterface interface {
	Create(ctx context.Context, req adminTimeSlotTemplateModel.CreateRequest, creatorID int64) (*adminTimeSlotTemplateModel.CreateResponse, error)
}

type DeleteServiceInterface interface {
	Delete(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.DeleteResponse, error)
}

type GetAllServiceInterface interface {
	GetAll(ctx context.Context, req adminTimeSlotTemplateModel.GetAllParsedRequest) (*adminTimeSlotTemplateModel.GetAllResponse, error)
}

type GetServiceInterface interface {
	Get(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.GetResponse, error)
}

type UpdateServiceInterface interface {
	Update(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateRequest) (*adminTimeSlotTemplateModel.UpdateResponse, error)
}
