package adminBooking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	db      *sqlx.DB
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, db *sqlx.DB) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID, bookingID int64, req adminBookingModel.UpdateParsedRequest) (*adminBookingModel.UpdateResponse, error) {
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	if !req.IsTimeSlotUpdateComplete() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingUpdateIncomplete)
	}

	// Get existing booking to verify it exists and is in SCHEDULED status
	existingBooking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	// Verify booking belongs to the store
	if existingBooking.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
	}
	// Verify booking is in SCHEDULED status
	if existingBooking.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	var oldTimeSlotID int64 = existingBooking.TimeSlotID
	var newServices []adminBookingModel.UpdateBookingServiceInfo

	if req.HasTimeSlotUpdate() {
		newServices, err = s.validateEntities(ctx, existingBooking.StylistID, existingBooking.TimeSlotID, *req.StylistID, *req.TimeSlotID, *req.MainServiceID, req.SubServiceIDs)
		if err != nil {
			return nil, err
		}
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Update booking record
	_, err = s.repo.Booking.UpdateBookingTx(ctx, tx, bookingID, sqlxRepo.UpdateBookingParams{
		StylistID:     req.StylistID,
		TimeSlotID:    req.TimeSlotID,
		IsChatEnabled: req.IsChatEnabled,
		Note:          req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking", err)
	}

	// Create new booking details if time slot is being updated
	if req.HasTimeSlotUpdate() {
		if err := s.updateBookingDetails(ctx, tx, bookingID, newServices); err != nil {
			return nil, err
		}
	}

	// Update time slot availability if changing time slot
	if req.TimeSlotID != nil && *req.TimeSlotID != oldTimeSlotID {
		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityByBookingIDTx(ctx, tx, oldTimeSlotID, true); err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to release old time slot", err)
		}

		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityByBookingIDTx(ctx, tx, *req.TimeSlotID, false); err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to reserve new time slot", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminBookingModel.UpdateResponse{
		ID: utils.FormatID(bookingID),
	}, nil
}

func (s *Update) validateEntities(ctx context.Context, oldStylistID, oldTimeSlotID int64, stylistID, timeSlotID, mainServiceID int64, subServiceIds []int64) ([]adminBookingModel.UpdateBookingServiceInfo, error) {
	// Validate stylist
	if oldStylistID != stylistID {
		_, err := s.queries.GetStylistByID(ctx, stylistID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
			}
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stylist", err)
		}
	}

	// Validate time slot
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get time slot", err)
	}
	if !timeSlot.IsAvailable.Bool && oldTimeSlotID != timeSlotID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingTimeSlotUnavailable)
	}

	// Validate main service
	mainService, err := s.queries.GetServiceByID(ctx, mainServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get main service", err)
	}
	if !mainService.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotActive)
	}
	if mainService.IsAddon.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotMainService)
	}

	// Validate sub services
	subServices := make([]dbgen.GetServiceByIdsRow, len(subServiceIds))
	if len(subServiceIds) > 0 {
		subServices, err = s.queries.GetServiceByIds(ctx, subServiceIds)
		if err != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		for i, subService := range subServices {
			if !subService.IsActive.Bool {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotActive)
			}
			if !subService.IsAddon.Bool {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotAddon)
			}
			subServices[i] = subService
		}
	}

	// if timeSlot time is not enough for service duration, return error
	endTime := utils.PgTimeToTime(timeSlot.EndTime)
	startTime := utils.PgTimeToTime(timeSlot.StartTime)

	timeSlotDuration := endTime.Sub(startTime)
	serviceDuration := time.Duration(mainService.DurationMinutes)
	for _, subService := range subServices {
		serviceDuration += time.Duration(subService.DurationMinutes)
	}

	if timeSlotDuration < serviceDuration {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotEnoughTime)
	}

	services := make([]adminBookingModel.UpdateBookingServiceInfo, len(subServices)+1)
	services[0] = adminBookingModel.UpdateBookingServiceInfo{
		ServiceId:     mainService.ID,
		ServiceName:   mainService.Name,
		IsMainService: true,
		Price:         mainService.Price,
	}
	for i, subService := range subServices {
		services[i+1] = adminBookingModel.UpdateBookingServiceInfo{
			ServiceId:     subService.ID,
			ServiceName:   subService.Name,
			IsMainService: false,
			Price:         subService.Price,
		}
	}

	return services, nil
}

func (s *Update) updateBookingDetails(ctx context.Context, tx *sqlx.Tx, bookingID int64, newServices []adminBookingModel.UpdateBookingServiceInfo) error {
	// Delete existing booking details
	if err := s.repo.BookingDetail.DeleteBookingDetailsByBookingIDTx(ctx, tx, bookingID); err != nil {
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete booking details", err)
	}

	// Create new booking details
	return s.createBookingDetails(ctx, tx, bookingID, newServices)
}

func (s *Update) createBookingDetails(ctx context.Context, tx *sqlx.Tx, bookingID int64, newServices []adminBookingModel.UpdateBookingServiceInfo) error {
	details := make([]sqlxRepo.BulkCreateBookingDetailsParams, len(newServices))

	for i, service := range newServices {
		detailID := utils.GenerateID()

		details[i] = sqlxRepo.BulkCreateBookingDetailsParams{
			ID:        detailID,
			BookingID: bookingID,
			ServiceID: service.ServiceId,
			Price:     service.Price,
		}
	}

	if err := s.repo.BookingDetail.BulkCreateBookingDetailsTx(ctx, tx, details); err != nil {
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create booking details", err)
	}

	return nil
}
