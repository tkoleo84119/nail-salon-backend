package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CancelBookingService struct {
	db   *sqlx.DB
	repo *sqlxRepo.Repositories
}

func NewCancelBookingService(db *sqlx.DB, repo *sqlxRepo.Repositories) CancelBookingServiceInterface {
	return &CancelBookingService{
		db:   db,
		repo: repo,
	}
}

func (s *CancelBookingService) CancelBooking(ctx context.Context, storeID, bookingID string, req adminBookingModel.CancelBookingRequest) (*adminBookingModel.CancelBookingResponse, error) {
	// Parse IDs
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid store ID", err)
	}

	bookingIDInt, err := utils.ParseID(bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid booking ID", err)
	}

	// Get existing booking to verify it exists and is in SCHEDULED status
	existingBooking, err := s.repo.Booking.GetByID(ctx, bookingIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	// Verify booking belongs to the store
	if existingBooking.StoreID != storeIDInt {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
	}
	// Verify booking is in SCHEDULED status
	if existingBooking.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToCancel)
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Cancel booking using repository
	id, err := s.repo.Booking.CancelBooking(ctx, tx, bookingIDInt, req.Status, req.CancelReason)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to cancel booking", err)
	}

	// Release time slot using repository
	err = s.repo.TimeSlot.UpdateTimeSlotAvailabilityTx(ctx, tx, existingBooking.TimeSlotID, true)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to release time slot", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminBookingModel.CancelBookingResponse{
		ID: utils.FormatID(id),
	}, nil
}
