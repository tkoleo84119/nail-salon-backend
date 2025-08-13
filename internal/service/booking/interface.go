package booking

import (
	"context"

	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
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

type CancelMyBookingServiceInterface interface {
	CancelMyBooking(ctx context.Context, bookingIDStr string, req bookingModel.CancelMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.CancelMyBookingResponse, error)
}

type GetMyBookingServiceInterface interface {
	GetMyBooking(ctx context.Context, bookingIDStr string, customerContext common.CustomerContext) (*bookingModel.GetMyBookingResponse, error)
}
