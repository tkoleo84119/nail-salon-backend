package adminSchedule

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetScheduleServiceInterface defines the interface for getting a single schedule
type GetScheduleServiceInterface interface {
	GetSchedule(ctx context.Context, storeID string, scheduleID string, staffContext common.StaffContext) (*adminScheduleModel.GetScheduleResponse, error)
}

type GetScheduleService struct {
	queries *dbgen.Queries
}

func NewGetScheduleService(queries *dbgen.Queries) *GetScheduleService {
	return &GetScheduleService{
		queries: queries,
	}
}

func (s *GetScheduleService) GetSchedule(ctx context.Context, storeID string, scheduleID string, staffContext common.StaffContext) (*adminScheduleModel.GetScheduleResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	// Parse schedule ID  
	scheduleIDInt, err := utils.ParseID(scheduleID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid schedule ID", err)
	}

	// Verify store exists
	_, err = s.queries.GetStoreByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Check store access for the staff member (except SUPER_ADMIN)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin {
		hasAccess, err := utils.CheckOneStoreAccess(storeIDInt, staffContext)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Get schedule by ID using SQLC (then validate store ID)
	schedule, err := s.queries.GetScheduleByID(ctx, scheduleIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule", err)
	}

	// Validate that the schedule belongs to the specified store
	if schedule.StoreID != storeIDInt {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	}

	// Get stylist information
	stylist, err := s.queries.GetStylistByID(ctx, schedule.StylistID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist", err)
	}

	// Get time slots for this schedule
	timeSlots, err := s.queries.GetTimeSlotsByScheduleID(ctx, scheduleIDInt)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get time slots", err)
	}

	// Build time slot response
	timeSlotItems := make([]adminScheduleModel.GetScheduleTimeSlotInfo, 0, len(timeSlots))
	for _, ts := range timeSlots {
		timeSlotItems = append(timeSlotItems, adminScheduleModel.GetScheduleTimeSlotInfo{
			ID:          utils.FormatID(ts.ID),
			StartTime:   utils.PgTimeToTimeString(ts.StartTime),
			EndTime:     utils.PgTimeToTimeString(ts.EndTime),
			IsAvailable: utils.PgBoolToBool(ts.IsAvailable),
		})
	}

	// Build response
	response := &adminScheduleModel.GetScheduleResponse{
		ID:       utils.FormatID(schedule.ID),
		WorkDate: utils.PgDateToDateString(schedule.WorkDate),
		Stylist: adminScheduleModel.GetScheduleStylistInfo{
			ID:   utils.FormatID(stylist.ID),
			Name: utils.PgTextToString(stylist.Name),
		},
		Note:      utils.PgTextToString(schedule.Note),
		TimeSlots: timeSlotItems,
	}

	return response, nil
}