package adminTimeSlotTemplate

import (
	"context"

	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateServiceInterface interface {
	Create(ctx context.Context, req adminTimeSlotTemplateModel.CreateRequest, creatorID int64) (*adminTimeSlotTemplateModel.CreateResponse, error)
}

type DeleteServiceInterface interface {
	Delete(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.DeleteResponse, error)
}

type UpdateTimeSlotTemplateServiceInterface interface {
	UpdateTimeSlotTemplate(ctx context.Context, templateID string, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error)
}

type CreateTimeSlotTemplateItemServiceInterface interface {
	CreateTimeSlotTemplateItem(ctx context.Context, templateID string, req adminTimeSlotTemplateModel.CreateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.CreateTimeSlotTemplateItemResponse, error)
}

type UpdateTimeSlotTemplateItemServiceInterface interface {
	UpdateTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateItemResponse, error)
}

type DeleteTimeSlotTemplateItemServiceInterface interface {
	DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.DeleteTimeSlotTemplateItemResponse, error)
}

type GetTimeSlotTemplateListServiceInterface interface {
	GetTimeSlotTemplateList(ctx context.Context, req adminTimeSlotTemplateModel.GetTimeSlotTemplateListParsedRequest) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse, error)
}

type GetTimeSlotTemplateServiceInterface interface {
	GetTimeSlotTemplate(ctx context.Context, templateID int64) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateResponse, error)
}
