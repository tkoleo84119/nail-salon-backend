package schedule

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
)

type CreateSchedulesBulkServiceInterface interface {
	CreateSchedulesBulk(ctx context.Context, req schedule.CreateSchedulesBulkRequest, staffContext common.StaffContext) (*schedule.CreateSchedulesBulkResponse, error)
}

type DeleteSchedulesBulkServiceInterface interface {
	DeleteSchedulesBulk(ctx context.Context, req schedule.DeleteSchedulesBulkRequest, staffContext common.StaffContext) (*schedule.DeleteSchedulesBulkResponse, error)
}

type CreateTimeSlotServiceInterface interface {
	CreateTimeSlot(ctx context.Context, scheduleID string, req schedule.CreateTimeSlotRequest, staffContext common.StaffContext) (*schedule.CreateTimeSlotResponse, error)
}

type UpdateTimeSlotServiceInterface interface {
	UpdateTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, req schedule.UpdateTimeSlotRequest, staffContext common.StaffContext) (*schedule.UpdateTimeSlotResponse, error)
}

type DeleteTimeSlotServiceInterface interface {
	DeleteTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, staffContext common.StaffContext) (*schedule.DeleteTimeSlotResponse, error)
}
