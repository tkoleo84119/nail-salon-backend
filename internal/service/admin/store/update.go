package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	repo *sqlxRepo.Repositories
}

func NewUpdate(repo *sqlxRepo.Repositories) *Update {
	return &Update{
		repo: repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID int64, req adminStoreModel.UpdateRequest, role string, storeIDs []int64) (*adminStoreModel.UpdateResponse, error) {
	// Validate role permissions (only SUPER_ADMIN and ADMIN can update stores)
	if role != adminStaffModel.RoleSuperAdmin && role != adminStaffModel.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Validate that at least one field has updates
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "at least one field must be provided for update", nil)
	}

	// Check if store exists
	_, err := s.repo.Store.GetStore(ctx, storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store details", err)
	}

	// For ADMIN role, check store access permission
	if role == adminStaffModel.RoleAdmin {
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
		nameExists, err := s.repo.Store.CheckStoreNameExistsExcluding(ctx, *req.Name, storeID)
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
