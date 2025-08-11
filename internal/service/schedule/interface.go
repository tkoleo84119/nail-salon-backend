package schedule

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, stylistID int64, req scheduleModel.GetAllParsedRequest, isBlacklisted bool) (*scheduleModel.GetAllResponse, error)
}

type GetTimeSlotServiceInterface interface {
	GetTimeSlotsBySchedule(ctx context.Context, scheduleID string, customerContext common.CustomerContext) (*scheduleModel.GetTimeSlotsByScheduleResponse, error)
}
