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
