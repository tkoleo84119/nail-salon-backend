package booking

import (
	"context"

	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateMyBookingServiceInterface interface {
	CreateMyBooking(ctx context.Context, req bookingModel.CreateMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.CreateMyBookingResponse, error)
}

type UpdateMyBookingServiceInterface interface {
	UpdateMyBooking(ctx context.Context, bookingIDStr string, req bookingModel.UpdateMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.UpdateMyBookingResponse, error)
}

type CancelMyBookingServiceInterface interface {
	CancelMyBooking(ctx context.Context, bookingIDStr string, req bookingModel.CancelMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.CancelMyBookingResponse, error)
}

type GetMyBookingsServiceInterface interface {
	GetMyBookings(ctx context.Context, queryParams bookingModel.GetMyBookingsQueryParams, customerContext common.CustomerContext) (*bookingModel.GetMyBookingsResponse, error)
}

type GetMyBookingServiceInterface interface {
	GetMyBooking(ctx context.Context, bookingIDStr string, customerContext common.CustomerContext) (*bookingModel.GetMyBookingResponse, error)
}
