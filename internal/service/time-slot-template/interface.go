package timeSlotTemplate

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
)

type CreateTimeSlotTemplateServiceInterface interface {
	CreateTimeSlotTemplate(ctx context.Context, req timeSlotTemplate.CreateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateResponse, error)
}

type UpdateTimeSlotTemplateServiceInterface interface {
	UpdateTimeSlotTemplate(ctx context.Context, templateID string, req timeSlotTemplate.UpdateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateResponse, error)
}

type CreateTimeSlotTemplateItemServiceInterface interface {
	CreateTimeSlotTemplateItem(ctx context.Context, templateID string, req timeSlotTemplate.CreateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateItemResponse, error)
}

type UpdateTimeSlotTemplateItemServiceInterface interface {
	UpdateTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, req timeSlotTemplate.UpdateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateItemResponse, error)
}

type DeleteTimeSlotTemplateItemServiceInterface interface {
	DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*timeSlotTemplate.DeleteTimeSlotTemplateItemResponse, error)
}
