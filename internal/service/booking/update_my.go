package booking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyBookingService struct {
	queries     dbgen.Querier
	bookingRepo *sqlxRepo.BookingRepository
	db          *pgxpool.Pool
}

func NewUpdateMyBookingService(queries dbgen.Querier, bookingRepo *sqlxRepo.BookingRepository, db *pgxpool.Pool) *UpdateMyBookingService {
	return &UpdateMyBookingService{
		queries:     queries,
		bookingRepo: bookingRepo,
		db:          db,
	}
}

func (s *UpdateMyBookingService) UpdateMyBooking(ctx context.Context, bookingIDStr string, req booking.UpdateMyBookingRequest, customerContext common.CustomerContext) (*booking.UpdateMyBookingResponse, error) {
	// Parse booking ID
	bookingID, err := utils.ParseID(bookingIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid booking ID", err)
	}

	// Validate that at least one field is provided
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	// Validate time slot update completeness
	if !req.IsTimeSlotUpdateComplete() {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "storeId, stylistId, timeSlotId, mainServiceId, and subServiceIds must be provided together", nil)
	}

	// Verify booking exists and belongs to customer
	bookingInfo, err := s.queries.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	if bookingInfo.CustomerID != customerContext.CustomerID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
	// only allow update booking in BookingStatusScheduled status
	if bookingInfo.Status != booking.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	// Parse IDs for validation
	var storeID, stylistID, timeSlotID, mainServiceID int64
	var subServiceIds []int64

	if req.HasTimeSlotUpdate() {
		storeID, err = utils.ParseID(*req.StoreId)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
		}

		stylistID, err = utils.ParseID(*req.StylistId)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid stylist ID", err)
		}

		timeSlotID, err = utils.ParseID(*req.TimeSlotId)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid time slot ID", err)
		}

		mainServiceID, err = utils.ParseID(*req.MainServiceId)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid main service ID", err)
		}

		if len(req.SubServiceIds) > 0 {
			subServiceIds, err = utils.ParseIDSlice(req.SubServiceIds)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid sub service IDs", err)
			}
		}

		if err := s.validateEntities(ctx, storeID, stylistID, timeSlotID, mainServiceID, subServiceIds); err != nil {
			return nil, err
		}
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Update booking
	_, err = s.bookingRepo.UpdateMyBooking(ctx, bookingID, customerContext.CustomerID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking", err)
	}

	// Update booking details if services are changing
	if req.HasTimeSlotUpdate() {
		if err := s.updateBookingDetails(ctx, qtx, bookingID, mainServiceID, subServiceIds, req); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Get updated booking with all details
	return s.buildResponse(ctx, bookingID)
}

func (s *UpdateMyBookingService) validateEntities(ctx context.Context, storeID, stylistID, timeSlotID, mainServiceID int64, subServiceIds []int64) error {
	// Validate store
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
	}
	if !store.IsActive.Bool {
		return errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	// Validate stylist
	_, err = s.queries.GetStylistByID(ctx, stylistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stylist", err)
	}

	// Validate time slot
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
		}
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get time slot", err)
	}
	if !timeSlot.IsAvailable.Bool {
		return errorCodes.NewServiceErrorWithCode(errorCodes.BookingTimeSlotUnavailable)
	}

	// Validate main service
	mainService, err := s.queries.GetServiceByID(ctx, mainServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get main service", err)
	}
	if !mainService.IsActive.Bool {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotActive)
	}

	// Validate sub services
	if len(subServiceIds) > 0 {
		subServices, err := s.queries.GetServiceByIds(ctx, subServiceIds)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get sub services", err)
		}
		for _, subService := range subServices {
			if !subService.IsActive.Bool {
				return errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotActive)
			}
		}
	}

	// if timeSlot time is not enough for service duration, return error
	endTime := utils.PgTimeToTime(timeSlot.EndTime)
	startTime := utils.PgTimeToTime(timeSlot.StartTime)
	timeSlotDuration := endTime.Sub(startTime)
	serviceDuration := time.Duration(mainService.DurationMinutes)
	if timeSlotDuration < serviceDuration {
		return errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotEnoughTime)
	}

	return nil
}

func (s *UpdateMyBookingService) updateBookingDetails(ctx context.Context, qtx *dbgen.Queries, bookingID, mainServiceID int64, subServiceIds []int64, req booking.UpdateMyBookingRequest) error {
	// Delete existing booking details
	if err := qtx.DeleteBookingDetailsByBookingID(ctx, bookingID); err != nil {
		return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete booking details", err)
	}

	// Prepare services for booking details
	var services []booking.BookingServiceInfo

	// Add main service
	if req.MainServiceId != nil {
		mainService, err := s.queries.GetServiceByID(ctx, mainServiceID)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get main service", err)
		}
		services = append(services, booking.BookingServiceInfo{
			ServiceId:     mainService.ID,
			ServiceName:   mainService.Name,
			Price:         utils.PgNumericToFloat64(mainService.Price),
			IsMainService: true,
		})
	}

	// Add sub services
	if len(subServiceIds) > 0 {
		subServices, err := s.queries.GetServiceByIds(ctx, subServiceIds)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get sub services", err)
		}
		for _, subService := range subServices {
			services = append(services, booking.BookingServiceInfo{
				ServiceId:   subService.ID,
				ServiceName: subService.Name,
				Price:       utils.PgNumericToFloat64(subService.Price),
			})
		}
	}

	// Create new booking details
	return s.createBookingDetails(ctx, qtx, bookingID, services)
}

func (s *UpdateMyBookingService) createBookingDetails(ctx context.Context, qtx *dbgen.Queries, bookingID int64, services []booking.BookingServiceInfo) error {
	for _, service := range services {
		detailID := utils.GenerateID()

		price, err := utils.Float64ToPgNumeric(service.Price)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert price", err)
		}

		_, err = qtx.CreateBookingDetail(ctx, dbgen.CreateBookingDetailParams{
			ID:        detailID,
			BookingID: bookingID,
			ServiceID: service.ServiceId,
			Price:     price,
		})
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create booking detail", err)
		}
	}
	return nil
}

func (s *UpdateMyBookingService) buildResponse(ctx context.Context, bookingID int64) (*booking.UpdateMyBookingResponse, error) {
	// Get complete booking info
	bookingInfo, err := s.queries.GetBookingByID(ctx, bookingID)
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
		if mainServiceName == "" {
			mainServiceName = detail.ServiceName
		} else {
			subServiceNames = append(subServiceNames, detail.ServiceName)
		}
	}

	return &booking.UpdateMyBookingResponse{
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
		IsChatEnabled:   bookingInfo.IsChatEnabled.Bool,
		Note:            &bookingInfo.Note.String,
		Status:          bookingInfo.Status,
		CreatedAt:       bookingInfo.CreatedAt.Time.String(),
		UpdatedAt:       bookingInfo.UpdatedAt.Time.String(),
	}, nil
}
