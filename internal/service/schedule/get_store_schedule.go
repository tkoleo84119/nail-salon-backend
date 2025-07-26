package schedule

import (
	"context"
	"database/sql"
	"errors"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreScheduleService struct {
	queries *dbgen.Queries
}

func NewGetStoreScheduleService(queries *dbgen.Queries) *GetStoreScheduleService {
	return &GetStoreScheduleService{
		queries: queries,
	}
}

func (s *GetStoreScheduleService) GetStoreSchedules(ctx context.Context, storeID, stylistID string, req scheduleModel.GetStoreSchedulesRequest, customerContext common.CustomerContext) (*scheduleModel.GetStoreSchedulesResponse, error) {
	// Input validation & ID parsing
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	stylistIDInt, err := utils.ParseID(stylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid stylist ID", err)
	}

	// Date validation
	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start date format", err)
	}

	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid end date format", err)
	}

	// date range validation (max 60 days)
	if endDate.Before(startDate) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValEndBeforeStart)
	}

	daysDiff := int(endDate.Sub(startDate).Hours() / 24)
	if daysDiff > 60 {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValDateRangeExceed60Days)
	}

	// Check if customer is blacklisted
	customer, err := s.queries.GetCustomerByID(ctx, customerContext.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer", err)
	}

	// If customer is blacklisted, return empty array
	if customer.IsBlacklisted.Bool {
		return &scheduleModel.GetStoreSchedulesResponse{
			Total: 0,
			Items: []scheduleModel.StoreScheduleResponseItem{},
		}, nil
	}

	// Data integrity validation - verify store exists and is active
	store, err := s.queries.GetStoreByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
	}

	if !store.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	// Verify stylist exists
	_, err = s.queries.GetStylistByID(ctx, stylistIDInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stylist", err)
	}

	// if startDate is in the past, move it to today
	if startDate.Before(time.Now()) {
		startDate = time.Now()
	}

	// Get available schedules for the store and stylist
	schedules, err := s.queries.GetAvailableSchedules(ctx, dbgen.GetAvailableSchedulesParams{
		StoreID:    storeIDInt,
		StylistID:  stylistIDInt,
		WorkDate:   utils.TimeToPgDate(startDate),
		WorkDate_2: utils.TimeToPgDate(endDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get schedules", err)
	}

	// Build response
	items := make([]scheduleModel.StoreScheduleResponseItem, 0, len(schedules))
	for _, schedule := range schedules {
		items = append(items, scheduleModel.StoreScheduleResponseItem{
			Date:           utils.PgDateToDateString(schedule.WorkDate),
			AvailableSlots: int(schedule.AvailableSlots),
		})
	}

	return &scheduleModel.GetStoreSchedulesResponse{
		Total: len(items),
		Items: items,
	}, nil
}
