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

type Cancel struct {
	db   *sqlx.DB
	repo *sqlxRepo.Repositories
}

func NewCancel(db *sqlx.DB, repo *sqlxRepo.Repositories) CancelInterface {
	return &Cancel{
		db:   db,
		repo: repo,
	}
}

func (s *Cancel) Cancel(ctx context.Context, storeID, bookingID int64, req adminBookingModel.CancelRequest) (*adminBookingModel.CancelResponse, error) {
	// Get existing booking to verify it exists and is in SCHEDULED status
	existingBooking, err := s.repo.Booking.GetByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	// Verify booking belongs to the store
	if existingBooking.StoreID != storeID {
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
	id, err := s.repo.Booking.CancelBooking(ctx, tx, bookingID, req.Status, req.CancelReason)
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

	return &adminBookingModel.CancelResponse{
		ID: utils.FormatID(id),
	}, nil
}
