package adminSchedule

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID int64, scheduleID int64, role string, storeIDs []int64) (*adminScheduleModel.GetResponse, error) {
	// Check store access for the staff member (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
			return nil, err
		}
	}

	// Get schedule by ID using SQLC (then validate store ID)
	rows, err := s.queries.GetScheduleWithTimeSlotsByID(ctx, scheduleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule", err)
	}

	response := adminScheduleModel.GetResponse{}
	response.TimeSlots = []adminScheduleModel.GetTimeSlotInfo{}
	for i, row := range rows {
		if i == 0 {
			response.ID = utils.FormatID(row.ID)
			response.WorkDate = utils.PgDateToDateString(row.WorkDate)
			response.Note = utils.PgTextToString(row.Note)
		}

		if row.TimeSlotID.Valid {
			response.TimeSlots = append(response.TimeSlots, adminScheduleModel.GetTimeSlotInfo{
				ID:          utils.FormatID(row.TimeSlotID.Int64),
				StartTime:   utils.PgTimeToTimeString(row.StartTime),
				EndTime:     utils.PgTimeToTimeString(row.EndTime),
				IsAvailable: utils.PgBoolToBool(row.IsAvailable),
			})
		}
	}

	return &response, nil
}
