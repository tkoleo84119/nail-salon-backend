package booking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateMyBookingService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreateMyBookingService(queries dbgen.Querier, db *pgxpool.Pool) *CreateMyBookingService {
	return &CreateMyBookingService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateMyBookingService) CreateMyBooking(ctx context.Context, req bookingModel.CreateMyBookingRequest, customerContext common.CustomerContext) (*bookingModel.CreateMyBookingResponse, error) {
	storeId, err := utils.ParseID(req.StoreId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid storeId", err)
	}

	stylistId, err := utils.ParseID(req.StylistId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid stylistId", err)
	}

	timeSlotId, err := utils.ParseID(req.TimeSlotId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid timeSlotId", err)
	}

	mainServiceId, err := utils.ParseID(req.MainServiceId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid mainServiceId", err)
	}

	var subServiceIds []int64
	if len(req.SubServiceIds) > 0 {
		subServiceIds, err = utils.ParseIDSlice(req.SubServiceIds)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid subServiceIds", err)
		}
	}

	// Check if store exists
	store, err := s.queries.GetStoreByID(ctx, storeId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StoreNotFound, "store not found", err)
	}
	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.StoreNotActive, "store is not active", err)
	}

	// Check if stylist exists
	stylist, err := s.queries.GetStylistByID(ctx, stylistId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.StylistNotFound, "stylist not found", err)
	}

	// Check if time slot exists
	timeSlot, err := s.queries.GetTimeSlotByID(ctx, timeSlotId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot not found", err)
	}
	if !timeSlot.IsAvailable.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.BookingTimeSlotUnavailable, "time slot is not available", err)
	}

	// check if schedule exists
	schedule, err := s.queries.GetScheduleByID(ctx, timeSlot.ScheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ScheduleNotFound, "schedule not found", err)
	}

	// Check if main service exists
	mainService, err := s.queries.GetServiceByID(ctx, mainServiceId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ServiceNotFound, "main service not found", err)
	}
	if !mainService.IsActive.Bool {
		return nil, errorCodes.NewServiceError(errorCodes.ServiceNotActive, "main service is not active", err)
	}

	// Check if sub services exist
	var subServices []dbgen.GetServiceByIdsRow
	if len(subServiceIds) > 0 {
		subServices, err = s.queries.GetServiceByIds(ctx, subServiceIds)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ServiceNotFound, "sub services not found", err)
		}
		for _, subService := range subServices {
			if !subService.IsActive.Bool {
				return nil, errorCodes.NewServiceError(errorCodes.ServiceNotActive, "sub service is not active", err)
			}
		}
	}

	// if timeSlot time is not enough for service duration, return error
	endTime := utils.PgTimeToTime(timeSlot.EndTime)
	startTime := utils.PgTimeToTime(timeSlot.StartTime)
	timeSlotDuration := endTime.Sub(startTime)
	serviceDuration := time.Duration(mainService.DurationMinutes)
	if timeSlotDuration < serviceDuration {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotEnoughTime)
	}

	services := make([]bookingModel.BookingServiceInfo, 0, len(subServices)+1)
	services = append(services, bookingModel.BookingServiceInfo{
		ServiceId:     mainService.ID,
		ServiceName:   mainService.Name,
		Price:         utils.PgNumericToFloat64(mainService.Price),
		IsMainService: true,
	})
	for _, subService := range subServices {
		services = append(services, bookingModel.BookingServiceInfo{
			ServiceId:   subService.ID,
			ServiceName: subService.Name,
			Price:       utils.PgNumericToFloat64(subService.Price),
		})
	}

	bookingId := utils.GenerateID()
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
		StoreID:       storeId,
		CustomerID:    customerContext.CustomerID,
		StylistID:     stylistId,
		TimeSlotID:    timeSlotId,
		IsChatEnabled: utils.BoolPtrToPgBool(req.IsChatEnabled),
		Note:          utils.StringPtrToPgText(req.Note, false),
		Status:        bookingModel.BookingStatusScheduled,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create booking failed", err)
	}

	// Create booking details
	err = s.createBookingDetails(ctx, qtx, bookingId, services)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "transaction commit failed", err)
	}

	subServiceNames := make([]string, 0, len(services)-1)
	for _, service := range services {
		if !service.IsMainService {
			subServiceNames = append(subServiceNames, service.ServiceName)
		}
	}

	response := &bookingModel.CreateMyBookingResponse{
		ID:              utils.FormatID(bookingId),
		StoreId:         utils.FormatID(storeId),
		StoreName:       store.Name,
		StylistId:       utils.FormatID(stylistId),
		StylistName:     utils.PgTextToString(stylist.Name),
		Date:            utils.PgDateToDateString(schedule.WorkDate),
		TimeSlotId:      utils.FormatID(timeSlotId),
		StartTime:       utils.PgTimeToTimeString(timeSlot.StartTime),
		EndTime:         utils.PgTimeToTimeString(timeSlot.EndTime),
		MainServiceName: services[0].ServiceName,
		SubServiceNames: subServiceNames,
		IsChatEnabled:   isChatEnabled,
		Note:            req.Note,
		Status:          bookingModel.BookingStatusScheduled,
		CreatedAt:       bookingInfo.CreatedAt.Time.String(),
		UpdatedAt:       bookingInfo.UpdatedAt.Time.String(),
	}

	return response, nil
}

func (s *CreateMyBookingService) createBookingDetails(ctx context.Context, qtx *dbgen.Queries, bookingId int64, services []bookingModel.BookingServiceInfo) error {
	for _, service := range services {
		detailId := utils.GenerateID()

		price, err := utils.Float64ToPgNumeric(service.Price)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert price", err)
		}

		_, err = qtx.CreateBookingDetail(ctx, dbgen.CreateBookingDetailParams{
			ID:        detailId,
			BookingID: bookingId,
			ServiceID: service.ServiceId,
			Price:     price,
		})
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.SysDatabaseError, "create booking detail failed", err)
		}
	}
	return nil
}
