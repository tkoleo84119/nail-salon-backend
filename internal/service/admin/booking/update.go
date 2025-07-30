package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateBookingByStaffService struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
	repo    *sqlx.Repositories
}

func NewUpdateBookingByStaffService(queries *dbgen.Queries, db *pgxpool.Pool, repo *sqlx.Repositories) UpdateBookingByStaffServiceInterface {
	return &UpdateBookingByStaffService{
		queries: queries,
		db:      db,
		repo:    repo,
	}
}

func (s *UpdateBookingByStaffService) UpdateBookingByStaff(ctx context.Context, storeID, bookingID string, req adminBookingModel.UpdateBookingByStaffRequest) (*adminBookingModel.UpdateBookingByStaffResponse, error) {
	// Parse IDs
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid store ID", err)
	}

	bookingIDInt, err := utils.ParseID(bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid booking ID", err)
	}

	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	if !req.IsTimeSlotUpdateComplete() {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "timeSlotId、mainServiceId、subServiceIds 必須一起傳入", nil)
	}

	// Get existing booking to verify it exists and is in SCHEDULED status
	existingBooking, err := s.queries.GetBookingByID(ctx, bookingIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}
	// Verify booking belongs to the store
	if existingBooking.StoreID != storeIDInt {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
	}
	// Verify booking is in SCHEDULED status
	if existingBooking.Status != bookingModel.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	var newTimeSlotIDInt *int64
	var oldTimeSlotID int64 = existingBooking.TimeSlotID
	// Handle time slot and service changes
	if req.HasTimeSlotUpdate() {
		// Validate time slot
		timeSlotIDInt, err := utils.ParseID(*req.TimeSlotID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid time slot ID", err)
		}
		newTimeSlotIDInt = &timeSlotIDInt

		// Get time slot details to verify it's available
		timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotIDInt)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get time slot", err)
		}

		// Check if time slot is available (only if it's different from current)
		if timeSlotIDInt != oldTimeSlotID && !timeSlot.IsAvailable.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingTimeSlotUnavailable)
		}

		// Validate main service
		mainServiceIDInt, err := utils.ParseID(*req.MainServiceID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid main service ID", err)
		}
		mainService, err := s.queries.GetServiceByID(ctx, mainServiceIDInt)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get main service", err)
		}
		if mainService.IsAddon.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotMainService)
		}

		// Validate sub services
		var subServiceIDInts []int64
		if len(req.SubServiceIDs) > 0 {
			for _, subServiceID := range req.SubServiceIDs {
				subServiceIDInt, err := utils.ParseID(subServiceID)
				if err != nil {
					return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid sub service ID", err)
				}
				subServiceIDInts = append(subServiceIDInts, subServiceIDInt)
			}

			// Validate all sub services exist and are addons
			services, err := s.queries.GetServiceByIds(ctx, subServiceIDInts)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get services", err)
			}

			if len(services) != len(subServiceIDInts) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
			}

			// Check that sub services are addons
			for _, service := range services {
				if !service.IsAddon.Bool {
					return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotAddon)
				}
			}
		}
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Update time slot availability if changing time slot
	if newTimeSlotIDInt != nil && *newTimeSlotIDInt != oldTimeSlotID {
		// Release old time slot
		oldTimeSlotParams := dbgen.UpdateTimeSlotParams{
			ID:          oldTimeSlotID,
			IsAvailable: utils.BoolToPgBool(true),
		}
		if _, err := qtx.UpdateTimeSlot(ctx, oldTimeSlotParams); err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to release old time slot", err)
		}

		// Reserve new time slot
		newTimeSlotParams := dbgen.UpdateTimeSlotParams{
			ID:          *newTimeSlotIDInt,
			IsAvailable: utils.BoolToPgBool(false),
		}
		if _, err := qtx.UpdateTimeSlot(ctx, newTimeSlotParams); err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to reserve new time slot", err)
		}
	}

	// Delete old booking details if time slot is being updated
	if req.HasTimeSlotUpdate() {
		if err := qtx.DeleteBookingDetailsByBookingID(ctx, bookingIDInt); err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete old booking details", err)
		}
	}

	// Create new booking details if time slot is being updated
	if req.HasTimeSlotUpdate() {
		mainServiceIDInt, _ := utils.ParseID(*req.MainServiceID)
		var subServiceIDInts []int64
		if len(req.SubServiceIDs) > 0 {
			for _, subServiceID := range req.SubServiceIDs {
				subServiceIDInt, _ := utils.ParseID(subServiceID)
				subServiceIDInts = append(subServiceIDInts, subServiceIDInt)
			}
		}

		allServiceIDs := append([]int64{mainServiceIDInt}, subServiceIDInts...)
		services, err := qtx.GetServiceByIds(ctx, allServiceIDs)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get services", err)
		}

		for _, service := range services {
			detailID := utils.GenerateID()

			createParams := dbgen.CreateBookingDetailParams{
				ID:        detailID,
				BookingID: bookingIDInt,
				ServiceID: service.ID,
				Price:     service.Price,
			}

			if _, err := qtx.CreateBookingDetail(ctx, createParams); err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create booking detail", err)
			}
		}
	}

	// Update booking record
	_, err = s.repo.Booking.UpdateBookingByStaff(ctx, bookingIDInt, storeIDInt, newTimeSlotIDInt, req.IsChatEnabled, req.Note)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminBookingModel.UpdateBookingByStaffResponse{
		ID: utils.FormatID(bookingIDInt),
	}, nil
}
