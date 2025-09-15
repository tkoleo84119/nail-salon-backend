package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateCompleted struct {
	queries *dbgen.Queries
}

func NewUpdateCompleted(queries *dbgen.Queries) UpdateCompletedInterface {
	return &UpdateCompleted{
		queries: queries,
	}
}

func (s *UpdateCompleted) UpdateCompleted(ctx context.Context, storeID, bookingID int64, req adminBookingModel.UpdateCompletedRequest, role string, updaterStoreIDs []int64) (*adminBookingModel.UpdateCompletedResponse, error) {
	if err := utils.CheckStoreAccess(storeID, updaterStoreIDs, role); err != nil {
		return nil, err
	}

	// Get existing booking to verify it exists and is in COMPLETED status
	existingBooking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
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

	// Verify booking is in COMPLETED status
	if existingBooking.Status != common.BookingStatusCompleted {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	// Update booking completed record
	err = s.queries.UpdateBookingActualDuration(ctx, dbgen.UpdateBookingActualDurationParams{
		ID:             bookingID,
		ActualDuration: utils.Int32PtrToPgInt4(req.ActualDuration),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update completed booking", err)
	}

	return &adminBookingModel.UpdateCompletedResponse{
		ID: utils.FormatID(bookingID),
	}, nil
}
