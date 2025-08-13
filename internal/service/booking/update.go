package booking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
	db      *sqlx.DB
}

func NewUpdate(queries dbgen.Querier, repo *sqlxRepo.Repositories, db *sqlx.DB) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
		db:      db,
	}
}

func (s *Update) Update(ctx context.Context, bookingID int64, req bookingModel.UpdateParsedRequest, customerID int64) (*bookingModel.UpdateResponse, error) {
	// Validate that at least one field is provided
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	// Validate time slot update completeness
	if !req.IsTimeSlotUpdateComplete() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingUpdateIncomplete)
	}

	// Verify booking exists and belongs to customer
	bookingInfo, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	if bookingInfo.CustomerID != customerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
	// only allow update booking in BookingStatusScheduled status
	if bookingInfo.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	var newServices []bookingModel.UpdateBookingServiceInfo
	if req.HasTimeSlotUpdate() {
		newServices, err = s.validateEntities(ctx, bookingInfo.StoreID, bookingInfo.StylistID, bookingInfo.TimeSlotID, *req.StoreId, *req.StylistId, *req.TimeSlotId, *req.MainServiceId, *req.SubServiceIds)
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

	// Update booking
	bookingID, err = s.repo.Booking.UpdateBookingTx(ctx, tx, bookingID, sqlxRepo.UpdateBookingParams{
		StoreID:       req.StoreId,
		StylistID:     req.StylistId,
		TimeSlotID:    req.TimeSlotId,
		IsChatEnabled: req.IsChatEnabled,
		Note:          req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking", err)
	}

	// Update booking details if services are changing
	if req.HasTimeSlotUpdate() {
		if err := s.updateBookingDetails(ctx, tx, bookingID, newServices); err != nil {
			return nil, err
		}
	}

	// when time slot is different, update old time slot to available and new time slot to unavailable
	if req.TimeSlotId != nil && bookingInfo.TimeSlotID != *req.TimeSlotId {
		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityByBookingIDTx(ctx, tx, bookingInfo.TimeSlotID, true); err != nil {
			return nil, err
		}
		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityByBookingIDTx(ctx, tx, *req.TimeSlotId, false); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Get updated booking with all details
	return s.buildResponse(ctx, bookingID)
}

func (s *Update) validateEntities(ctx context.Context, oldStoreID, oldStylistID, oldTimeSlotID int64, storeID, stylistID, timeSlotID, mainServiceID int64, subServiceIds []int64) ([]bookingModel.UpdateBookingServiceInfo, error) {
	// Validate store
	if oldStoreID != storeID {
		store, err := s.queries.GetStoreByID(ctx, storeID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
			}
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
		}
		if !store.IsActive.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
		}
	}

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

	services := make([]bookingModel.UpdateBookingServiceInfo, len(subServices)+1)
	services[0] = bookingModel.UpdateBookingServiceInfo{
		ServiceId:     mainService.ID,
		ServiceName:   mainService.Name,
		IsMainService: true,
		Price:         mainService.Price,
	}
	for i, subService := range subServices {
		services[i+1] = bookingModel.UpdateBookingServiceInfo{
			ServiceId:     subService.ID,
			ServiceName:   subService.Name,
			IsMainService: false,
			Price:         subService.Price,
		}
	}

	return services, nil
}

func (s *Update) updateBookingDetails(ctx context.Context, tx *sqlx.Tx, bookingID int64, newServices []bookingModel.UpdateBookingServiceInfo) error {
	// Delete existing booking details
	if err := s.repo.BookingDetail.DeleteBookingDetailsByBookingIDTx(ctx, tx, bookingID); err != nil {
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete booking details", err)
	}

	// Create new booking details
	return s.createBookingDetails(ctx, tx, bookingID, newServices)
}

func (s *Update) createBookingDetails(ctx context.Context, tx *sqlx.Tx, bookingID int64, newServices []bookingModel.UpdateBookingServiceInfo) error {
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

func (s *Update) buildResponse(ctx context.Context, bookingID int64) (*bookingModel.UpdateResponse, error) {
	// Get complete booking info
	bookingInfo, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get updated booking", err)
	}

	// Get booking details for services
	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking details", err)
	}

	// Separate main and sub services
	var mainServiceName string
	var subServiceNames []string

	// Assuming first service is main service (you might need better logic here)
	for _, detail := range bookingDetails {
		if detail.IsAddon.Bool {
			subServiceNames = append(subServiceNames, detail.ServiceName)
		} else {
			mainServiceName = detail.ServiceName
		}
	}

	return &bookingModel.UpdateResponse{
		ID:              utils.FormatID(bookingInfo.ID),
		StoreId:         utils.FormatID(bookingInfo.StoreID),
		StoreName:       bookingInfo.StoreName,
		StylistId:       utils.FormatID(bookingInfo.StylistID),
		StylistName:     utils.PgTextToString(bookingInfo.StylistName),
		Date:            utils.PgDateToDateString(bookingInfo.WorkDate),
		TimeSlotId:      utils.FormatID(bookingInfo.TimeSlotID),
		StartTime:       utils.PgTimeToTimeString(bookingInfo.StartTime),
		EndTime:         utils.PgTimeToTimeString(bookingInfo.EndTime),
		MainServiceName: mainServiceName,
		SubServiceNames: subServiceNames,
		IsChatEnabled:   utils.PgBoolToBool(bookingInfo.IsChatEnabled),
		Note:            utils.PgTextToString(bookingInfo.Note),
		Status:          bookingInfo.Status,
		CreatedAt:       utils.PgTimestamptzToTimeString(bookingInfo.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(bookingInfo.UpdatedAt),
	}, nil
}
