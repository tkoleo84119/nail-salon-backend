package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries dbgen.Querier, repo *sqlxRepo.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID int64, req adminStoreModel.UpdateRequest, role string, storeIDs []int64) (*adminStoreModel.UpdateResponse, error) {
	// Validate role permissions (only SUPER_ADMIN and ADMIN can update stores)
	if role != common.RoleSuperAdmin && role != common.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Validate that at least one field has updates
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// For ADMIN role, check store access permission
	if role == common.RoleAdmin {
		hasAccess, err := utils.CheckStoreAccess(storeID, storeIDs)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
		}
		if !hasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if name is unique (excluding current store)
	if req.Name != nil {
		nameExists, err := s.queries.CheckStoreNameExistsExcluding(ctx, dbgen.CheckStoreNameExistsExcludingParams{
			ID:   storeID,
			Name: *req.Name,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store name uniqueness", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreAlreadyExists)
		}
	}

	// Update store using sqlx repository
	updatedStore, err := s.repo.Store.UpdateStore(ctx, storeID, sqlxRepo.UpdateStoreParams{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update store", err)
	}

	response := &adminStoreModel.UpdateResponse{
		ID:        utils.FormatID(updatedStore.ID),
		Name:      updatedStore.Name,
		Address:   utils.PgTextToString(updatedStore.Address),
		Phone:     utils.PgTextToString(updatedStore.Phone),
		IsActive:  utils.PgBoolToBool(updatedStore.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(updatedStore.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(updatedStore.UpdatedAt),
	}

	return response, nil
}
