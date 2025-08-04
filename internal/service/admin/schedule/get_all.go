package adminSchedule

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetScheduleListService struct {
	repo *sqlxRepo.Repositories
}

func NewGetScheduleListService(repo *sqlxRepo.Repositories) *GetScheduleListService {
	return &GetScheduleListService{
		repo: repo,
	}
}

func (s *GetScheduleListService) GetScheduleList(ctx context.Context, storeID int64, req adminScheduleModel.GetScheduleListParsedRequest, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleListResponse, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleEndBeforeStart)
	}

	if req.EndDate.Sub(req.StartDate) > 31*24*time.Hour {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleDateRangeExceed31Days)
	}

	// Check store exists
	_, err := s.repo.Store.GetStoreByID(ctx, storeID, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for the staff member (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		hasAccess, err := utils.CheckStoreAccess(storeID, storeIDs)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Get schedules from repository with dynamic filtering
	rows, err := s.repo.Schedule.GetStoreScheduleByDateRange(ctx, storeID, req.StartDate, req.EndDate, sqlxRepo.GetStoreScheduleByDateRangeParams{
		StylistID:   req.StylistID,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule list", err)
	}

	stylistMap := make(map[int64]*adminScheduleModel.GetScheduleListStylistItem)
	scheduleMap := make(map[int64]map[int64]*adminScheduleModel.GetScheduleListScheduleItem)

	for _, row := range rows {
		stylistID := row.StylistID

		stylist, exists := stylistMap[stylistID]
		if !exists {
			stylist = &adminScheduleModel.GetScheduleListStylistItem{
				ID:        utils.FormatID(stylistID),
				Name:      utils.PgTextToString(row.StylistName),
				Schedules: []adminScheduleModel.GetScheduleListScheduleItem{},
			}

			stylistMap[stylistID] = stylist
			scheduleMap[stylistID] = make(map[int64]*adminScheduleModel.GetScheduleListScheduleItem)
		}

		workID := row.ID
		sMap := scheduleMap[stylistID]
		schedule, ok := sMap[workID]
		if !ok {
			schedule = &adminScheduleModel.GetScheduleListScheduleItem{
				ID:        utils.FormatID(workID),
				WorkDate:  utils.PgDateToDateString(row.WorkDate),
				Note:      utils.PgTextToString(row.Note),
				TimeSlots: []adminScheduleModel.GetScheduleListTimeSlotInfo{},
			}
			stylist.Schedules = append(stylist.Schedules, *schedule)
			// update scheduleMap with the new schedule pointer
			schedule = &stylist.Schedules[len(stylist.Schedules)-1]
			sMap[workID] = schedule
		}

		schedule.TimeSlots = append(schedule.TimeSlots, adminScheduleModel.GetScheduleListTimeSlotInfo{
			ID:          utils.FormatID(row.TimeSlotID),
			StartTime:   utils.PgTimeToTimeString(row.StartTime),
			EndTime:     utils.PgTimeToTimeString(row.EndTime),
			IsAvailable: utils.PgBoolToBool(row.IsAvailable),
		})
	}

	response := adminScheduleModel.GetScheduleListResponse{
		StylistList: make([]adminScheduleModel.GetScheduleListStylistItem, 0, len(stylistMap)),
	}

	for _, stylist := range stylistMap {
		response.StylistList = append(response.StylistList, *stylist)
	}

	return &response, nil
}
