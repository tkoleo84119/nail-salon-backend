package booking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Cancel struct {
	queries dbgen.Querier
}

func NewCancel(queries dbgen.Querier) CancelInterface {
	return &Cancel{
		queries: queries,
	}
}

func (s *Cancel) Cancel(ctx context.Context, bookingID int64, req bookingModel.CancelRequest, customerID int64) (*bookingModel.CancelResponse, error) {
	// Verify booking exists and belongs to customer
	bookingInfo, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Check if booking belongs to the customer
	if bookingInfo.CustomerID != customerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if booking is in a cancelable state (only SCHEDULED bookings can be canceled)
	if bookingInfo.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToCancel)
	}

	// Cancel booking with optional cancel reason
	_, err = s.queries.CancelBooking(ctx, dbgen.CancelBookingParams{
		ID:           bookingID,
		Status:       bookingModel.BookingStatusCancelled,
		CancelReason: utils.StringPtrToPgText(req.CancelReason, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to cancel booking", err)
	}

	newBooking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Build response
	return &bookingModel.CancelResponse{
		ID:          utils.FormatID(newBooking.ID),
		StoreId:     utils.FormatID(newBooking.StoreID),
		StoreName:   newBooking.StoreName,
		StylistId:   utils.FormatID(newBooking.StylistID),
		StylistName: utils.PgTextToString(newBooking.StylistName),
		Date:        utils.PgDateToDateString(newBooking.WorkDate),
		TimeSlotId:  utils.FormatID(newBooking.TimeSlotID),
		StartTime:   utils.PgTimeToTimeString(newBooking.StartTime),
		EndTime:     utils.PgTimeToTimeString(newBooking.EndTime),
		Status:      newBooking.Status,
		CreatedAt:   utils.PgTimestamptzToTimeString(newBooking.CreatedAt),
		UpdatedAt:   utils.PgTimestamptzToTimeString(newBooking.UpdatedAt),
	}, nil
}
