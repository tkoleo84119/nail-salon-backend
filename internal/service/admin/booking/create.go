package adminBooking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
}

func NewCreate(queries *dbgen.Queries, db *pgxpool.Pool) *Create {
	return &Create{
		queries: queries,
		db:      db,
	}
}

func (s *Create) Create(ctx context.Context, storeID int64, req adminBookingModel.CreateParsedRequest, role string, storeIds []int64) (*adminBookingModel.CreateResponse, error) {
	// Verify store exists
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	// Check store access for staff (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIds); err != nil {
			return nil, err
		}
	}

	// Verify customer exists
	exists, err := s.queries.CheckCustomerExistsByID(ctx, req.CustomerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check customer exists", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
	}

	// Verify time slot exists, is available, and belongs to the stylist
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, req.TimeSlotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot", err)
	}
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.BookingTimeSlotUnavailable, "Time slot is not available", nil)
	}

	// Get the schedule to obtain the date for the time slot
	schedule, err := s.queries.GetScheduleByID(ctx, timeSlot.ScheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule for time slot", err)
	}
	// Verify schedule belongs to the store and stylist
	if schedule.StoreID != storeID {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotBelongToStore, "Time slot does not belong to the specified store", nil)
	}

	// Verify main service exists (services are global, not store-specific)
	mainService, err := s.queries.GetServiceByID(ctx, req.MainServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get main service", err)
	}

	// Verify sub services exist (services are global, not store-specific)
	var subServices []dbgen.Service
	for _, subServiceIDInt := range req.SubServiceIDs {
		subService, err := s.queries.GetServiceByID(ctx, subServiceIDInt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
			}
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get sub service", err)
		}

		subServices = append(subServices, subService)
	}

	services := make([]bookingModel.CreateBookingServiceInfo, len(subServices)+1)
	services[0] = bookingModel.CreateBookingServiceInfo{
		ServiceId:     mainService.ID,
		ServiceName:   mainService.Name,
		IsMainService: true,
		Price:         mainService.Price,
	}
	for i, subService := range subServices {
		services[i+1] = bookingModel.CreateBookingServiceInfo{
			ServiceId:     subService.ID,
			ServiceName:   subService.Name,
			IsMainService: false,
			Price:         subService.Price,
		}
	}

	bookingId := utils.GenerateID()
	bookingDetails, err := s.parseBookingDetails(bookingId, services)
	if err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Create booking
	booking, err := qtx.CreateBooking(ctx, dbgen.CreateBookingParams{
		ID:            bookingId,
		StoreID:       storeID,
		CustomerID:    req.CustomerID,
		StylistID:     schedule.StylistID,
		TimeSlotID:    req.TimeSlotID,
		IsChatEnabled: utils.BoolPtrToPgBool(&req.IsChatEnabled),
		Note:          utils.StringPtrToPgText(req.Note, true),
		Status:        common.BookingStatusScheduled,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to create booking", err)
	}

	// Create booking details for main service
	_, err = qtx.CreateBookingDetails(ctx, bookingDetails)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to create main booking detail", err)
	}

	isAvailable := false
	// Mark time slot as unavailable
	_, err = qtx.UpdateTimeSlotIsAvailable(ctx, dbgen.UpdateTimeSlotIsAvailableParams{
		ID:          req.TimeSlotID,
		IsAvailable: utils.BoolPtrToPgBool(&isAvailable),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to update time slot availability", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to commit transaction", err)
	}

	// Build response
	response := &adminBookingModel.CreateResponse{
		ID: utils.FormatID(booking.ID),
	}

	return response, nil
}

func (s *Create) parseBookingDetails(bookingId int64, services []bookingModel.CreateBookingServiceInfo) ([]dbgen.CreateBookingDetailsParams, error) {
	bookingDetails := make([]dbgen.CreateBookingDetailsParams, len(services))
	now := time.Now()
	nowPg := utils.TimePtrToPgTimestamptz(&now)

	for i, service := range services {
		bookingDetails[i] = dbgen.CreateBookingDetailsParams{
			ID:        utils.GenerateID(),
			BookingID: bookingId,
			ServiceID: service.ServiceId,
			Price:     service.Price,
			CreatedAt: nowPg,
			UpdatedAt: nowPg,
		}
	}
	return bookingDetails, nil
}
