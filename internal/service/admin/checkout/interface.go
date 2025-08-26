package adminCheckout

import (
	"context"

	adminCheckoutModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/checkout"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, bookingID int64, req adminCheckoutModel.CreateParsedRequest, creatorID int64, storeIDs []int64) (*adminCheckoutModel.CreateResponse, error)
}
