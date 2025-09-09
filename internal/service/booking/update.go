package booking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries       *dbgen.Queries
	repo          *sqlxRepo.Repositories
	db            *sqlx.DB
	lineMessenger *utils.LineMessageClient
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, db *sqlx.DB, lineMessenger *utils.LineMessageClient) *Update {
	return &Update{
		queries:       queries,
		repo:          repo,
		db:            db,
		lineMessenger: lineMessenger,
	}
}

func (s *Update) Update(ctx context.Context, bookingID int64, req bookingModel.UpdateParsedRequest, customerID int64) (*bookingModel.UpdateResponse, error) {
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
	if bookingInfo.Status != common.BookingStatusScheduled {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	var newServices []bookingModel.UpdateBookingServiceInfo
	if req.HasTimeSlotUpdate() {
		newServices, err = s.validateEntities(ctx, bookingInfo.StoreID, bookingInfo.StylistID, bookingInfo.TimeSlotID, *req.TimeSlotId, *req.MainServiceId, *req.SubServiceIds)
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
	bookingID, err = s.repo.Booking.UpdateBookingTx(ctx, tx, bookingID, sqlxRepo.UpdateBookingTxParams{
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
		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityTx(ctx, tx, bookingInfo.TimeSlotID, true); err != nil {
			return nil, err
		}
		if err := s.repo.TimeSlot.UpdateTimeSlotAvailabilityTx(ctx, tx, *req.TimeSlotId, false); err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// if customer not have chat permission (this mean customer not give permission to liff app, so can't send message in liff app) and update time slot, send line message
	needSendLineMessage := req.HasChatPermission != nil && !*req.HasChatPermission && req.HasTimeSlotUpdate()

	// Get updated booking with all details
	return s.buildResponse(ctx, bookingID, needSendLineMessage)
}

func (s *Update) validateEntities(ctx context.Context, oldStoreID, oldStylistID, oldTimeSlotID int64, timeSlotID, mainServiceID int64, subServiceIds []int64) ([]bookingModel.UpdateBookingServiceInfo, error) {
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
	endTime, err := utils.PgTimeToTime(timeSlot.EndTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert time", err)
	}
	startTime, err := utils.PgTimeToTime(timeSlot.StartTime)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert time", err)
	}

	timeSlotDuration := endTime.Sub(startTime)
	serviceDuration := time.Duration(mainService.DurationMinutes) * time.Minute
	for _, subService := range subServices {
		serviceDuration += time.Duration(subService.DurationMinutes) * time.Minute
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

func (s *Update) buildResponse(ctx context.Context, bookingID int64, needSendLineMessage bool) (*bookingModel.UpdateResponse, error) {
	// get new booking info, because booking info may be changed by update booking
	bookingInfo, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking details", err)
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

	response := &bookingModel.UpdateResponse{
		ID:              utils.FormatID(bookingInfo.ID),
		StoreId:         utils.FormatID(bookingInfo.StoreID),
		StoreName:       bookingInfo.StoreName,
		StylistId:       utils.FormatID(bookingInfo.StylistID),
		StylistName:     utils.PgTextToString(bookingInfo.StylistName),
		CustomerName:    bookingInfo.CustomerName,
		CustomerPhone:   bookingInfo.CustomerPhone,
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
	}

	if needSendLineMessage {
		err := s.lineMessenger.SendBookingNotification(bookingInfo.CustomerLineUid, common.BookingActionUpdated, &utils.BookingData{
			StoreName:       response.StoreName,
			Date:            response.Date,
			StartTime:       response.StartTime,
			EndTime:         response.EndTime,
			CustomerName:    &response.CustomerName,
			CustomerPhone:   &response.CustomerPhone,
			StylistName:     response.StylistName,
			MainServiceName: response.MainServiceName,
			SubServiceNames: response.SubServiceNames,
		})
		if err != nil {
			log.Printf("failed to send line message: %v", err)
		}
	}

	return response, nil
}
