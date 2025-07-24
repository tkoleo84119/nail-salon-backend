package booking

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateMyBookingServiceInterface interface {
	CreateMyBooking(ctx context.Context, req booking.CreateMyBookingRequest, customerContext common.CustomerContext) (*booking.CreateMyBookingResponse, error)
}

type UpdateMyBookingServiceInterface interface {
	UpdateMyBooking(ctx context.Context, bookingIDStr string, req booking.UpdateMyBookingRequest, customerContext common.CustomerContext) (*booking.UpdateMyBookingResponse, error)
}

type CancelMyBookingServiceInterface interface {
	CancelMyBooking(ctx context.Context, bookingIDStr string, req booking.CancelMyBookingRequest, customerContext common.CustomerContext) (*booking.CancelMyBookingResponse, error)
}
