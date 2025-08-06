package adminTimeSlot

import (
	"context"

	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
)

type CreateInterface interface {
	Create(ctx context.Context, scheduleID int64, req adminTimeSlotModel.CreateRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminTimeSlotModel.CreateResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, scheduleID int64, timeSlotID int64, req adminTimeSlotModel.UpdateRequest, updaterID int64, updaterRole string, updaterStoreIDs []int64) (*adminTimeSlotModel.UpdateResponse, error)
}
