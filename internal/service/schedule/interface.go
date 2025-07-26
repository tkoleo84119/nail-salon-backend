package schedule

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
)

type ScheduleServiceInterface interface {
	GetStoreSchedules(ctx context.Context, storeID, stylistID string, req scheduleModel.GetStoreSchedulesRequest, customerContext common.CustomerContext) (*scheduleModel.GetStoreSchedulesResponse, error)
}
