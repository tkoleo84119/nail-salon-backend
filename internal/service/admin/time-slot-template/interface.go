package adminTimeSlotTemplate

import (
	"context"

	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateTimeSlotTemplateServiceInterface interface {
	CreateTimeSlotTemplate(ctx context.Context, req adminTimeSlotTemplateModel.CreateTimeSlotTemplateRequest, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.CreateTimeSlotTemplateResponse, error)
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

type DeleteTimeSlotTemplateServiceInterface interface {
	DeleteTimeSlotTemplate(ctx context.Context, templateID string, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.DeleteTimeSlotTemplateResponse, error)
}

type DeleteTimeSlotTemplateItemServiceInterface interface {
	DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*adminTimeSlotTemplateModel.DeleteTimeSlotTemplateItemResponse, error)
}

type GetTimeSlotTemplateListServiceInterface interface {
	GetTimeSlotTemplateList(ctx context.Context, req adminTimeSlotTemplateModel.GetTimeSlotTemplateListParsedRequest) (*adminTimeSlotTemplateModel.GetTimeSlotTemplateListResponse, error)
}
