package booking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CancelMyBookingService struct {
	queries dbgen.Querier
}

func NewCancelMyBookingService(queries dbgen.Querier) CancelMyBookingServiceInterface {
	return &CancelMyBookingService{
		queries: queries,
	}
}

func (s *CancelMyBookingService) CancelMyBooking(ctx context.Context, bookingIDStr string, req bookingModel.CancelMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.CancelMyBookingResponse, error) {
	// Parse booking ID
	bookingID, err := utils.ParseID(bookingIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid booking ID", err)
	}

	// Verify booking exists and belongs to customer
	bookingInfo, err := s.queries.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Check if booking belongs to the customer
	if bookingInfo.CustomerID != customerContext.CustomerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if booking is in a cancelable state (only SCHEDULED bookings can be canceled)
	if bookingInfo.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToCancel)
	}

	// Cancel booking with optional cancel reason
	cancelledBooking, err := s.queries.CancelBooking(ctx, dbgen.CancelBookingParams{
		ID:           bookingID,
		Status:       bookingModel.BookingStatusCancelled,
		CancelReason: utils.StringPtrToPgText(req.CancelReason, true), // Empty as NULL
		CustomerID:   customerContext.CustomerID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to cancel booking", err)
	}

	// Build response
	var cancelReason *string
	if cancelledBooking.CancelReason.Valid {
		cancelReason = &cancelledBooking.CancelReason.String
	}

	return &bookingModel.CancelMyBookingResponse{
		ID:           utils.FormatID(cancelledBooking.ID),
		Status:       cancelledBooking.Status,
		CancelReason: cancelReason,
	}, nil
}
