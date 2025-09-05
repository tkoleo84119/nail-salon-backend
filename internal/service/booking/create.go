package booking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries       *dbgen.Queries
	db            *pgxpool.Pool
	lineMessenger *utils.LineMessageClient
}

func NewCreate(queries *dbgen.Queries, db *pgxpool.Pool, lineConfig config.LineConfig) *Create {
	lineMessenger := utils.NewLineMessenger(lineConfig.MessageAccessToken)
	return &Create{
		queries:       queries,
		db:            db,
		lineMessenger: lineMessenger,
	}
}

func (s *Create) Create(ctx context.Context, req bookingModel.CreateParsedRequest, customerID int64) (*bookingModel.CreateResponse, error) {
	// Check if store exists
	store, err := s.queries.GetStoreByID(ctx, req.StoreId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "get store by id failed", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	// Check if customer exists
	customer, err := s.queries.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
	}
	if customer.IsBlacklisted.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerIsBlacklisted)
	}

	// Check if stylist exists
	stylistName, err := s.queries.GetActiveStylistNameByID(ctx, req.StylistId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "get stylist by id failed", err)
	}

	// Check if time slot exists
	timeSlot, err := s.queries.GetTimeSlotWithScheduleByID(ctx, req.TimeSlotId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "get time slot by id failed", err)
	}

	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingTimeSlotUnavailable)
	}

	// Check if main service exists
	mainService, err := s.queries.GetServiceByID(ctx, req.MainServiceId)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
	}
	if !mainService.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotActive)
	}
	if mainService.IsAddon.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotMainService)
	}

	// Check if sub services exist
	subServices := make([]dbgen.GetServiceByIdsRow, len(req.SubServiceIds))
	if len(req.SubServiceIds) > 0 {
		subServices, err = s.queries.GetServiceByIds(ctx, req.SubServiceIds)
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
	isChatEnabled := req.IsChatEnabled != nil && *req.IsChatEnabled

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction begin failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Create booking
	bookingInfo, err := qtx.CreateBooking(ctx, dbgen.CreateBookingParams{
		ID:            bookingId,
		StoreID:       req.StoreId,
		CustomerID:    customerID,
		StylistID:     req.StylistId,
		TimeSlotID:    req.TimeSlotId,
		IsChatEnabled: utils.BoolPtrToPgBool(&isChatEnabled),
		Note:          utils.StringPtrToPgText(req.Note, true),
		Status:        common.BookingStatusScheduled,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create booking failed", err)
	}

	// Create booking details
	_, err = qtx.CreateBookingDetails(ctx, bookingDetails)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create booking details failed", err)
	}

	// update time slot to unavailable
	isAvailable := false
	_, err = qtx.UpdateTimeSlotIsAvailable(ctx, dbgen.UpdateTimeSlotIsAvailableParams{
		ID:          req.TimeSlotId,
		IsAvailable: utils.BoolPtrToPgBool(&isAvailable),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "update time slot failed", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction commit failed", err)
	}

	subServiceNames := make([]string, len(subServices))
	for i, service := range subServices {
		if service.IsAddon.Bool {
			subServiceNames[i] = service.Name
		}
	}

	response := &bookingModel.CreateResponse{
		ID:              utils.FormatID(bookingId),
		StoreId:         utils.FormatID(req.StoreId),
		StoreName:       store.Name,
		StylistId:       utils.FormatID(req.StylistId),
		StylistName:     utils.PgTextToString(stylistName),
		CustomerName:    customer.Name,
		CustomerPhone:   customer.Phone,
		Date:            utils.PgDateToDateString(timeSlot.WorkDate),
		TimeSlotId:      utils.FormatID(req.TimeSlotId),
		StartTime:       utils.PgTimeToTimeString(timeSlot.StartTime),
		EndTime:         utils.PgTimeToTimeString(timeSlot.EndTime),
		MainServiceName: services[0].ServiceName,
		SubServiceNames: subServiceNames,
		IsChatEnabled:   isChatEnabled,
		Note:            utils.PgTextToString(bookingInfo.Note),
		Status:          common.BookingStatusScheduled,
		CreatedAt:       utils.PgTimestamptzToTimeString(bookingInfo.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(bookingInfo.UpdatedAt),
	}

	// if customer no chat permission (this mean customer not give permission to liff app, so can't send message in liff app) send line message, but not return error
	if req.HasChatPermission != nil && !*req.HasChatPermission {
		err := s.lineMessenger.SendBookingNotification(customer.LineUid, common.BookingActionCreated, &utils.BookingData{
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

func (s *Create) parseBookingDetails(bookingId int64, services []bookingModel.CreateBookingServiceInfo) ([]dbgen.CreateBookingDetailsParams, error) {
	bookingDetails := make([]dbgen.CreateBookingDetailsParams, len(services))
	now := utils.TimeToPgTimestamptz(time.Now())

	for i, service := range services {
		bookingDetails[i] = dbgen.CreateBookingDetailsParams{
			ID:        utils.GenerateID(),
			BookingID: bookingId,
			ServiceID: service.ServiceId,
			Price:     service.Price,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return bookingDetails, nil
}
