package adminSchedule

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetScheduleService struct {
	repo *sqlxRepo.Repositories
}

func NewGetScheduleService(repo *sqlxRepo.Repositories) *GetScheduleService {
	return &GetScheduleService{
		repo: repo,
	}
}

func (s *GetScheduleService) GetSchedule(ctx context.Context, storeID int64, scheduleID int64, role string, storeIDs []int64) (*adminScheduleModel.GetScheduleResponse, error) {
	// Verify store exists
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

	// Get schedule by ID using SQLC (then validate store ID)
	rows, err := s.repo.Schedule.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule", err)
	}

	response := adminScheduleModel.GetScheduleResponse{}
	response.TimeSlots = []adminScheduleModel.GetScheduleTimeSlotInfo{}
	for i, row := range rows {
		if i == 0 {
			response.ID = utils.FormatID(row.ID)
			response.WorkDate = utils.PgDateToDateString(row.WorkDate)
			response.Note = utils.PgTextToString(row.Note)
		}

		if row.TimeSlotID.Valid {
			response.TimeSlots = append(response.TimeSlots, adminScheduleModel.GetScheduleTimeSlotInfo{
				ID:          utils.FormatID(row.TimeSlotID.Int64),
				StartTime:   utils.PgTimeToTimeString(row.StartTime),
				EndTime:     utils.PgTimeToTimeString(row.EndTime),
				IsAvailable: utils.PgBoolToBool(row.IsAvailable),
			})
		}
	}

	return &response, nil
}
