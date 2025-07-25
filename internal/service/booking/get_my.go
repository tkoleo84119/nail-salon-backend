package booking

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyBookingsService struct {
	bookingRepo *sqlxRepo.BookingRepository
}

func NewGetMyBookingsService(bookingRepo *sqlxRepo.BookingRepository) GetMyBookingsServiceInterface {
	return &GetMyBookingsService{
		bookingRepo: bookingRepo,
	}
}

func (s *GetMyBookingsService) GetMyBookings(ctx context.Context, queryParams bookingModel.GetMyBookingsQueryParams, customerContext common.CustomerContext) (*bookingModel.GetMyBookingsResponse, error) {
	if len(queryParams.Status) > 0 {
		for _, status := range queryParams.Status {
			if !bookingModel.IsValidBookingStatus(status) {
				return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid status value", nil)
			}
		}
	}

	// Set default status if none provided (according to API doc, default is SCHEDULED)
	statuses := queryParams.Status
	if len(statuses) == 0 {
		statuses = []string{bookingModel.BookingStatusScheduled}
	}

	// Get bookings from repository
	bookings, total, err := s.bookingRepo.GetMyBookings(ctx, customerContext.CustomerID, statuses, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get bookings", err)
	}

	// Build response items
	items := make([]bookingModel.GetMyBookingsItemModel, len(bookings))
	for i, booking := range bookings {
		items[i] = bookingModel.GetMyBookingsItemModel{
			ID:          utils.FormatID(booking.ID),
			StoreId:     utils.FormatID(booking.StoreID),
			StoreName:   booking.StoreName,
			StylistId:   utils.FormatID(booking.StylistID),
			StylistName: booking.StylistName,
			Date:        booking.Date,
			TimeSlot: bookingModel.GetMyBookingsTimeSlotModel{
				ID:        utils.FormatID(booking.TimeSlotID),
				StartTime: booking.StartTime,
				EndTime:   booking.EndTime,
			},
			Status: booking.Status,
		}
	}

	return &bookingModel.GetMyBookingsResponse{
		Total: total,
		Items: items,
	}, nil
}
