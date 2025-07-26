package schedule

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
)

type GetStoreScheduleServiceInterface interface {
	GetStoreSchedules(ctx context.Context, storeID, stylistID string, req scheduleModel.GetStoreSchedulesRequest, customerContext common.CustomerContext) (*scheduleModel.GetStoreSchedulesResponse, error)
}

type GetTimeSlotServiceInterface interface {
	GetTimeSlotsBySchedule(ctx context.Context, scheduleID string, customerContext common.CustomerContext) (*scheduleModel.GetTimeSlotsByScheduleResponse, error)
}
