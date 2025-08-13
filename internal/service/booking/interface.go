package booking

import (
	"context"

	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
)

type CreateInterface interface {
	Create(ctx context.Context, req bookingModel.CreateParsedRequest, customerID int64) (*bookingModel.CreateResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, bookingID int64, req bookingModel.UpdateParsedRequest, customerID int64) (*bookingModel.UpdateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, queryParams bookingModel.GetAllParsedRequest, customerID int64) (*bookingModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, bookingID int64, customerID int64) (*bookingModel.GetResponse, error)
}

type CancelInterface interface {
	Cancel(ctx context.Context, bookingID int64, req bookingModel.CancelRequest, customerID int64) (*bookingModel.CancelResponse, error)
}
