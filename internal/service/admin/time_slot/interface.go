package adminTimeSlot

import (
	"context"

	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
)

type CreateInterface interface {
	Create(ctx context.Context, scheduleID int64, req adminTimeSlotModel.CreateRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminTimeSlotModel.CreateResponse, error)
}
