package adminService

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetServiceService struct {
	queries *dbgen.Queries
}

func NewGetServiceService(queries *dbgen.Queries) *GetServiceService {
	return &GetServiceService{
		queries: queries,
	}
}

func (s *GetServiceService) GetService(ctx context.Context, storeID, serviceID string, staffContext common.StaffContext) (*adminServiceModel.GetServiceResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid store ID", err)
	}

	// Parse service ID
	serviceIDInt, err := utils.ParseID(serviceID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "Invalid service ID", err)
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

	// Get service information
	service, err := s.queries.GetServiceByID(ctx, serviceIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get service", err)
	}

	// Build response
	response := &adminServiceModel.GetServiceResponse{
		ID:              utils.FormatID(service.ID),
		Name:            service.Name,
		DurationMinutes: service.DurationMinutes,
		Price:           int64(utils.PgNumericToFloat64(service.Price)),
		IsAddon:         service.IsAddon.Bool,
		IsActive:        service.IsActive.Bool,
		IsVisible:       service.IsVisible.Bool,
		Note:            service.Note.String,
	}

	return response, nil
}
