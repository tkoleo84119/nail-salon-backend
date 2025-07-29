package adminBooking

import (
	"context"

	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
)

type UpdateBookingByStaffServiceInterface interface {
	UpdateBookingByStaff(ctx context.Context, storeID, bookingID string, req adminBookingModel.UpdateBookingByStaffRequest) (*adminBookingModel.UpdateBookingByStaffResponse, error)
}

type CancelBookingServiceInterface interface {
	CancelBooking(ctx context.Context, storeID, bookingID string, req adminBookingModel.CancelBookingRequest) (*adminBookingModel.CancelBookingResponse, error)
}
