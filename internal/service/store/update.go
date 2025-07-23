package store

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStoreService struct {
	queries         dbgen.Querier
	storeRepository sqlxRepo.StoreRepositoryInterface
}

func NewUpdateStoreService(queries dbgen.Querier, storeRepository sqlxRepo.StoreRepositoryInterface) *UpdateStoreService {
	return &UpdateStoreService{
		queries:         queries,
		storeRepository: storeRepository,
	}
}

func (s *UpdateStoreService) UpdateStore(ctx context.Context, storeID string, req store.UpdateStoreRequest, staffContext common.StaffContext) (*store.UpdateStoreResponse, error) {
	// Parse store ID
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	// Validate that at least one field has updates
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Validate role permissions (only SUPER_ADMIN and ADMIN can update stores)
	if staffContext.Role != staff.RoleSuperAdmin && staffContext.Role != staff.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store exists
	_, err = s.queries.GetStoreDetailByID(ctx, parsedStoreID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store details", err)
	}

	// For ADMIN role, check store access permission
	if staffContext.Role == staff.RoleAdmin {
		storeIDs := []int64{}
		for _, store := range staffContext.StoreList {
			parsedStoreID, err := utils.ParseID(store.ID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
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
	response, err := s.storeRepository.UpdateStore(ctx, parsedStoreID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update store", err)
	}

	return response, nil
}
