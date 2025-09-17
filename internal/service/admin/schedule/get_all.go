package adminSchedule

import (
	"context"
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAllService struct {
	repo    *sqlxRepo.Repositories
	queries *dbgen.Queries
}

func NewGetAll(queries *dbgen.Queries, repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAllService{
		repo:    repo,
		queries: queries,
	}
}

func (s *GetAllService) GetAll(ctx context.Context, storeID int64, req adminScheduleModel.GetAllParsedRequest, role string, storeIDs []int64) (*adminScheduleModel.GetAllResponse, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleEndBeforeStart)
	}

	if req.EndDate.Sub(req.StartDate) > 31*24*time.Hour {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleDateRangeExceed31Days)
	}

	// Check store exists
	exists, err := s.queries.CheckStoreExistByID(ctx, storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check store exists", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	// Check store access for the staff member
	if err := utils.CheckStoreAccess(storeID, storeIDs, role); err != nil {
		return nil, err
	}

	// Get schedules from repository with dynamic filtering
	rows, err := s.repo.Schedule.GetStoreSchedulesByDateRange(ctx, storeID, req.StartDate, req.EndDate, sqlxRepo.GetStoreSchedulesByDateRangeParams{
		StylistID:   req.StylistID,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule list", err)
	}

	stylistMap := make(map[int64]*adminScheduleModel.GetAllStylistItem)
	scheduleMap := make(map[int64]map[int64]*adminScheduleModel.GetAllScheduleItem)

	for _, row := range rows {
		stylistID := row.StylistID

		stylist, exists := stylistMap[stylistID]
		if !exists {
			stylist = &adminScheduleModel.GetAllStylistItem{
				ID:        utils.FormatID(stylistID),
				Name:      utils.PgTextToString(row.StylistName),
				Schedules: []adminScheduleModel.GetAllScheduleItem{},
			}

			stylistMap[stylistID] = stylist
			scheduleMap[stylistID] = make(map[int64]*adminScheduleModel.GetAllScheduleItem)
		}

		workID := row.ID
		sMap := scheduleMap[stylistID]
		schedule, ok := sMap[workID]
		if !ok {
			schedule = &adminScheduleModel.GetAllScheduleItem{
				ID:        utils.FormatID(workID),
				WorkDate:  utils.PgDateToDateString(row.WorkDate),
				Note:      utils.PgTextToString(row.Note),
				TimeSlots: []adminScheduleModel.GetAllTimeSlotInfo{},
			}
			stylist.Schedules = append(stylist.Schedules, *schedule)
			// update scheduleMap with the new schedule pointer
			schedule = &stylist.Schedules[len(stylist.Schedules)-1]
			sMap[workID] = schedule
		}

		if row.TimeSlotID != nil {
			schedule.TimeSlots = append(schedule.TimeSlots, adminScheduleModel.GetAllTimeSlotInfo{
				ID:          utils.FormatID(*row.TimeSlotID),
				StartTime:   utils.PgTimeToTimeString(row.StartTime),
				EndTime:     utils.PgTimeToTimeString(row.EndTime),
				IsAvailable: utils.PgBoolToBool(row.IsAvailable),
			})
		}
	}

	response := adminScheduleModel.GetAllResponse{
		StylistList: make([]adminScheduleModel.GetAllStylistItem, 0, len(stylistMap)),
	}

	for _, stylist := range stylistMap {
		response.StylistList = append(response.StylistList, *stylist)
	}

	return &response, nil
}
