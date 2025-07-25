package adminSchedule

import (
	"context"

	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateSchedulesBulkServiceInterface interface {
	CreateSchedulesBulk(ctx context.Context, req adminScheduleModel.CreateSchedulesBulkRequest, staffContext common.StaffContext) (*adminScheduleModel.CreateSchedulesBulkResponse, error)
}

type DeleteSchedulesBulkServiceInterface interface {
	DeleteSchedulesBulk(ctx context.Context, req adminScheduleModel.DeleteSchedulesBulkRequest, staffContext common.StaffContext) (*adminScheduleModel.DeleteSchedulesBulkResponse, error)
}

type CreateTimeSlotServiceInterface interface {
	CreateTimeSlot(ctx context.Context, scheduleID string, req adminScheduleModel.CreateTimeSlotRequest, staffContext common.StaffContext) (*adminScheduleModel.CreateTimeSlotResponse, error)
}

type UpdateTimeSlotServiceInterface interface {
	UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req adminScheduleModel.UpdateTimeSlotRequest, staffContext common.StaffContext) (*adminScheduleModel.UpdateTimeSlotResponse, error)
}

type DeleteTimeSlotServiceInterface interface {
	DeleteTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, staffContext common.StaffContext) (*adminScheduleModel.DeleteTimeSlotResponse, error)
}
