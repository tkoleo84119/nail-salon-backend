package adminTimeSlotTemplateItem

import (
	"context"

	adminTimeSlotTemplateItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot_template_item"
)

type CreateServiceInterface interface {
	Create(ctx context.Context, templateID int64, req adminTimeSlotTemplateItemModel.CreateRequest) (*adminTimeSlotTemplateItemModel.CreateResponse, error)
}

type UpdateServiceInterface interface {
	Update(ctx context.Context, templateID int64, itemID int64, req adminTimeSlotTemplateItemModel.UpdateRequest) (*adminTimeSlotTemplateItemModel.UpdateResponse, error)
}
