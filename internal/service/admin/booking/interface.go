package adminBooking

import (
	"context"

	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminBookingModel.CreateParsedRequest, role string, storeIds []int64) (*adminBookingModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminBookingModel.GetAllParsedRequest, role string, storeIds []int64) (*adminBookingModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, storeID, bookingID int64, role string, storeIds []int64) (*adminBookingModel.GetResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, storeID, bookingID int64, req adminBookingModel.UpdateParsedRequest) (*adminBookingModel.UpdateResponse, error)
}

type CancelInterface interface {
	Cancel(ctx context.Context, storeID, bookingID int64, req adminBookingModel.CancelRequest) (*adminBookingModel.CancelResponse, error)
}
