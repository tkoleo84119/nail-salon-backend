package adminSchedule

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// GetScheduleListServiceInterface defines the interface for getting schedule list
type GetScheduleListServiceInterface interface {
	GetScheduleList(ctx context.Context, storeID string, req adminScheduleModel.GetScheduleListRequest, staffContext common.StaffContext) (*adminScheduleModel.GetScheduleListResponse, error)
}

type GetScheduleListService struct {
	queries      *dbgen.Queries
	scheduleRepo sqlxRepo.ScheduleRepositoryInterface
}

func NewGetScheduleListService(queries *dbgen.Queries, scheduleRepo sqlxRepo.ScheduleRepositoryInterface) *GetScheduleListService {
	return &GetScheduleListService{
		queries:      queries,
		scheduleRepo: scheduleRepo,
	}
}

func (s *GetScheduleListService) GetScheduleList(ctx context.Context, storeID string, req adminScheduleModel.GetScheduleListRequest, staffContext common.StaffContext) (*adminScheduleModel.GetScheduleListResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	// Parse stylist ID
	var stylistIDInt int64
	if req.StylistID != nil {
		stylistIDInt, err = utils.ParseID(*req.StylistID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid stylist ID", err)
		}
	}

	// validate startDate and endDate
	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid start date", err)
	}
	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid end date", err)
	}

	if endDate.Before(startDate) {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "End date must be after start date", nil)
	}

	if endDate.Sub(startDate) > 60*24*time.Hour {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Date range cannot exceed 60 days", nil)
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

	// Get schedules from repository with dynamic filtering
	response, err := s.scheduleRepo.GetStoreScheduleList(ctx, storeIDInt, sqlxRepo.GetStoreScheduleListParams{
		StylistID:   &stylistIDInt,
		StartDate:   startDate,
		EndDate:     endDate,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get schedule list", err)
	}

	return response, nil
}
