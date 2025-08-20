package adminBooking

import (
	"context"

	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminBookingModel.GetAllParsedRequest, role string, storeIds []int64) (*adminBookingModel.GetAllResponse, error)
}

type UpdateBookingByStaffServiceInterface interface {
	UpdateBookingByStaff(ctx context.Context, storeID, bookingID string, req adminBookingModel.UpdateBookingByStaffRequest) (*adminBookingModel.UpdateBookingByStaffResponse, error)
}

type CancelBookingServiceInterface interface {
	CancelBooking(ctx context.Context, storeID, bookingID string, req adminBookingModel.CancelBookingRequest) (*adminBookingModel.CancelBookingResponse, error)
}
