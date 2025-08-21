package booking

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, queryParams bookingModel.GetAllParsedRequest, customerID int64) (*bookingModel.GetAllResponse, error) {
	// Get bookings from repository
	total, bookings, err := s.repo.Booking.GetAllCustomerBookingsByFilter(ctx, customerID, sqlxRepo.GetAllCustomerBookingsByFilterParams{
		Limit:  &queryParams.Limit,
		Offset: &queryParams.Offset,
		Sort:   &queryParams.Sort,
		Status: *queryParams.Status,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get bookings", err)
	}

	// Build response items
	items := make([]bookingModel.GetAllItem, len(bookings))
	for i, booking := range bookings {
		items[i] = bookingModel.GetAllItem{
			ID:          utils.FormatID(booking.ID),
			StoreId:     utils.FormatID(booking.StoreID),
			StoreName:   booking.StoreName,
			StylistId:   utils.FormatID(booking.StylistID),
			StylistName: booking.StylistName,
			Date:        utils.PgDateToDateString(booking.Date),
			TimeSlotId:  utils.FormatID(booking.TimeSlotID),
			StartTime:   utils.PgTimeToTimeString(booking.StartTime),
			EndTime:     utils.PgTimeToTimeString(booking.EndTime),
			Status:      booking.Status,
		}
	}

	return &bookingModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
