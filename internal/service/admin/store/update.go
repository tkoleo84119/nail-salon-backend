package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStoreService struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
}

func NewUpdateStoreService(queries dbgen.Querier, repo *sqlxRepo.Repositories) *UpdateStoreService {
	return &UpdateStoreService{
		queries: queries,
		repo:    repo,
	}
}

func (s *UpdateStoreService) UpdateStore(ctx context.Context, storeID string, req adminStoreModel.UpdateStoreRequest, staffContext common.StaffContext) (*adminStoreModel.UpdateStoreResponse, error) {
	// Parse store ID
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid store ID", err)
	}

	// Validate that at least one field has updates
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Validate role permissions (only SUPER_ADMIN and ADMIN can update stores)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin && staffContext.Role != adminStaffModel.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store exists
	_, err = s.queries.GetStoreDetailByID(ctx, parsedStoreID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store details", err)
	}

	// For ADMIN role, check store access permission
	if staffContext.Role == adminStaffModel.RoleAdmin {
		storeIDs := []int64{}
		for _, store := range staffContext.StoreList {
			parsedStoreID, err := utils.ParseID(store.ID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid store ID", err)
			}
			storeIDs = append(storeIDs, parsedStoreID)
		}

		hasAccess := false
		for _, storeID := range storeIDs {
			if storeID == parsedStoreID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if name is unique (excluding current store)
	if req.Name != nil {
		nameExists, err := s.queries.CheckStoreNameExistsExcluding(ctx, dbgen.CheckStoreNameExistsExcludingParams{
			Name: *req.Name,
			ID:   parsedStoreID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store name uniqueness", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreAlreadyExists)
		}
	}

	// Update store using sqlx repository
	response, err := s.repo.Store.UpdateStore(ctx, parsedStoreID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update store", err)
	}

	return response, nil
}
