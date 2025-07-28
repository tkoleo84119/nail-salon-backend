package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// CreateBookingServiceInterface defines the interface for creating bookings
type CreateBookingServiceInterface interface {
	CreateBooking(ctx context.Context, storeID string, req adminBookingModel.CreateBookingRequest, staffContext common.StaffContext) (*adminBookingModel.CreateBookingResponse, error)
}

type CreateBookingService struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
}

func NewCreateBookingService(queries *dbgen.Queries, db *pgxpool.Pool) *CreateBookingService {
	return &CreateBookingService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateBookingService) CreateBooking(ctx context.Context, storeID string, req adminBookingModel.CreateBookingRequest, staffContext common.StaffContext) (*adminBookingModel.CreateBookingResponse, error) {
	// Parse and validate IDs
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	customerIDInt, err := utils.ParseID(req.CustomerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid customer ID", err)
	}

	timeSlotIDInt, err := utils.ParseID(req.TimeSlotID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid time slot ID", err)
	}

	mainServiceIDInt, err := utils.ParseID(req.MainServiceID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid main service ID", err)
	}

	// Parse sub service IDs
	var subServiceIDInts []int64
	for _, subServiceID := range req.SubServiceIDs {
		subServiceIDInt, err := utils.ParseID(subServiceID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid sub service ID", err)
		}
		subServiceIDInts = append(subServiceIDInts, subServiceIDInt)
	}

	// Verify store exists
	store, err := s.queries.GetStoreByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for staff (except SUPER_ADMIN)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin {
		hasAccess, err := utils.CheckOneStoreAccess(storeIDInt, staffContext)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Verify customer exists
	_, err = s.queries.GetCustomerByID(ctx, customerIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get customer", err)
	}

	// Verify time slot exists, is available, and belongs to the stylist
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slot", err)
	}
	if !utils.PgBoolToBool(timeSlot.IsAvailable) {
		return nil, errorCodes.NewServiceError(errorCodes.BookingTimeSlotUnavailable, "Time slot is not available", nil)
	}

	// Get the schedule to obtain the date for the time slot
	schedule, err := s.queries.GetScheduleByID(ctx, timeSlot.ScheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule for time slot", err)
	}
	// Verify schedule belongs to the store and stylist
	if schedule.StoreID != storeIDInt {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotBelongToStore, "Time slot does not belong to the specified store", nil)
	}

	// Verify main service exists (services are global, not store-specific)
	mainService, err := s.queries.GetServiceByID(ctx, mainServiceIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get main service", err)
	}

	// Verify sub services exist (services are global, not store-specific)
	var subServices []dbgen.GetServiceByIDRow
	for _, subServiceIDInt := range subServiceIDInts {
		subService, err := s.queries.GetServiceByID(ctx, subServiceIDInt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
			}
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get sub service", err)
		}

		subServices = append(subServices, subService)
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
		ID:            utils.GenerateID(),
		StoreID:       storeIDInt,
		CustomerID:    customerIDInt,
		StylistID:     schedule.StylistID,
		TimeSlotID:    timeSlotIDInt,
		IsChatEnabled: utils.BoolToPgBool(req.IsChatEnabled),
		Note:          utils.StringPtrToPgText(req.Note, false),
		Status:        bookingModel.BookingStatusScheduled,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to create booking", err)
	}

	// Create booking details for main service
	_, err = qtx.CreateBookingDetail(ctx, dbgen.CreateBookingDetailParams{
		ID:        utils.GenerateID(),
		BookingID: booking.ID,
		ServiceID: mainServiceIDInt,
		Price:     mainService.Price,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to create main booking detail", err)
	}

	// Create booking details for sub services
	for i, subService := range subServices {
		_, err := qtx.CreateBookingDetail(ctx, dbgen.CreateBookingDetailParams{
			ID:        utils.GenerateID(),
			BookingID: booking.ID,
			ServiceID: subServiceIDInts[i],
			Price:     subService.Price,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to create sub booking detail", err)
		}
	}

	// Mark time slot as unavailable
	_, err = qtx.UpdateTimeSlot(ctx, dbgen.UpdateTimeSlotParams{
		ID:          timeSlotIDInt,
		IsAvailable: utils.BoolToPgBool(false),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to update time slot availability", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to commit transaction", err)
	}

	// Build sub service names
	subServiceNames := make([]string, 0, len(subServices))
	for _, subService := range subServices {
		subServiceNames = append(subServiceNames, subService.Name)
	}

	stylist, err := s.queries.GetStylistByID(ctx, schedule.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist", err)
	}

	// Build response
	response := &adminBookingModel.CreateBookingResponse{
		ID:              utils.FormatID(booking.ID),
		StoreID:         utils.FormatID(store.ID),
		StoreName:       store.Name,
		StylistID:       utils.FormatID(schedule.StylistID),
		StylistName:     utils.PgTextToString(stylist.Name),
		Date:            utils.PgDateToDateString(schedule.WorkDate),
		TimeSlotID:      utils.FormatID(timeSlot.ID),
		StartTime:       utils.PgTimeToTimeString(timeSlot.StartTime),
		EndTime:         utils.PgTimeToTimeString(timeSlot.EndTime),
		MainServiceName: mainService.Name,
		SubServiceNames: subServiceNames,
		IsChatEnabled:   req.IsChatEnabled,
		Note:            utils.PgTextToString(booking.Note),
		Status:          booking.Status,
		CreatedAt:       booking.CreatedAt.Time.String(),
		UpdatedAt:       booking.UpdatedAt.Time.String(),
	}

	return response, nil
}
