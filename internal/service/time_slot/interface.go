package timeSlot

import (
	"context"

	timeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/time_slot"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, scheduleID int64, isBlacklisted bool) (*timeSlotModel.GetAllResponse, error)
}
