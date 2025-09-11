package schedule

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
}

func NewGetAll(queries *dbgen.Queries) GetAllInterface {
	return &GetAll{
		queries: queries,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, stylistID int64, req scheduleModel.GetAllParsedRequest, isBlacklisted bool) (*scheduleModel.GetAllResponse, error) {
	// date range validation (max 60 days)
	if req.EndDate.Before(req.StartDate) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleEndBeforeStart)
	}

	daysDiff := int(req.EndDate.Sub(req.StartDate).Hours() / 24)
	if daysDiff > 31 {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleDateRangeExceed31Days)
	}

	exists, err := s.queries.CheckStoreExistAndActive(ctx, storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store exist and active", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	// Verify stylist exists
	exist, err := s.queries.CheckStylistExistAndActive(ctx, stylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stylist exist and active", err)
	}
	if !exist {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	}

	schedules := make([]scheduleModel.GetAllItem, 0)
	response := scheduleModel.GetAllResponse{
		Schedules: schedules,
	}

	// If customer is blacklisted, return empty array
	if isBlacklisted {
		return &response, nil
	}

	// if startDate is in the past, move it to today
	startDate := req.StartDate
	if startDate.Before(time.Now()) {
		startDate = time.Now()
	}

	rows, err := s.queries.GetAvailableSchedules(ctx, dbgen.GetAvailableSchedulesParams{
		StoreID:    storeID,
		StylistID:  stylistID,
		WorkDate:   utils.TimePtrToPgDate(&startDate),
		WorkDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get available schedules", err)
	}

	for _, row := range rows {
		schedules = append(schedules, scheduleModel.GetAllItem{
			ID:   utils.FormatID(row.ID),
			Date: utils.PgDateToDateString(row.WorkDate),
		})
	}

	response.Schedules = schedules

	return &response, nil
}
