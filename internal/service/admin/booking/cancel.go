package adminBooking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Cancel struct {
	queries     *dbgen.Queries
	db          *sqlx.DB
	repo        *sqlxRepo.Repositories
	activityLog cache.ActivityLogCacheInterface
}

func NewCancel(
	queries *dbgen.Queries,
	db *sqlx.DB,
	repo *sqlxRepo.Repositories,
	activityLog cache.ActivityLogCacheInterface,
) CancelInterface {
	return &Cancel{
		queries:     queries,
		db:          db,
		repo:        repo,
		activityLog: activityLog,
	}
}

func (s *Cancel) Cancel(ctx context.Context, storeID, bookingID int64, req adminBookingModel.CancelRequest, staffName string) (*adminBookingModel.CancelResponse, error) {
	// Get existing booking to verify it exists and is in SCHEDULED status
	booking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	// Verify booking belongs to the store
	if booking.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
	}
	// Verify booking is in SCHEDULED status
	if booking.Status != common.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToCancel)
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Cancel booking using repository
	id, err := s.repo.Booking.CancelBookingTx(ctx, tx, bookingID, req.Status, req.CancelReason)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to cancel booking", err)
	}

	// Release time slot using repository
	err = s.repo.TimeSlot.UpdateTimeSlotAvailabilityTx(ctx, tx, booking.TimeSlotID, true)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to release time slot", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Log activity
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.activityLog.LogAdminBookingCancel(logCtx, staffName, booking.CustomerName, utils.PgTextToString(booking.CustomerLineName), booking.StoreName); err != nil {
			log.Printf("failed to log admin booking cancel activity: %v", err)
		}
	}()

	return &adminBookingModel.CancelResponse{
		ID: utils.FormatID(id),
	}, nil
}
