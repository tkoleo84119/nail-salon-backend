package booking

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateMyBookingServiceInterface interface {
	CreateMyBooking(ctx context.Context, req booking.CreateMyBookingRequest, customerContext common.CustomerContext) (*booking.CreateMyBookingResponse, error)
}
