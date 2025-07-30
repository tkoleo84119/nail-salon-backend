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
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetServiceListService struct {
	queries *dbgen.Queries
	repo    *sqlx.Repositories
}

func NewGetServiceListService(queries *dbgen.Queries, repo *sqlx.Repositories) *GetServiceListService {
	return &GetServiceListService{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetServiceListService) GetServiceList(ctx context.Context, storeID string, req adminServiceModel.GetServiceListRequest, staffContext common.StaffContext) (*adminServiceModel.GetServiceListResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid store ID", err)
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
		var storeAccess []int64
		for _, store := range staffContext.StoreList {
			storeId, err := utils.ParseID(store.ID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid store ID", err)
			}
			storeAccess = append(storeAccess, storeId)
		}

		hasAccess := false
		for _, storeId := range storeAccess {
			if storeId == storeIDInt {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// set default value
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	// Get service list from repository
	results, total, err := s.repo.Service.GetStoreServiceList(ctx, storeIDInt, sqlx.GetStoreServiceListParams{
		Name:      req.Name,
		IsAddon:   req.IsAddon,
		IsActive:  req.IsActive,
		IsVisible: req.IsVisible,
		Limit:     &limit,
		Offset:    &offset,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get service list", err)
	}

	items := make([]adminServiceModel.ServiceListItemDTO, len(results))
	for i, result := range results {
		items[i] = adminServiceModel.ServiceListItemDTO{
			ID:              utils.FormatID(result.ID),
			Name:            result.Name,
			Price:           result.Price,
			DurationMinutes: result.DurationMinutes,
			IsAddon:         result.IsAddon,
			IsActive:        result.IsActive,
			IsVisible:       result.IsVisible,
			Note:            result.Note,
		}
	}

	response := &adminServiceModel.GetServiceListResponse{
		Total: total,
		Items: items,
	}

	return response, nil
}
