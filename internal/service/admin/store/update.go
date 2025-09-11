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
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID int64, req adminStoreModel.UpdateRequest, role string, storeIDs []int64) (*adminStoreModel.UpdateResponse, error) {
	// For ADMIN role, check store access permission
	if role == common.RoleAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
			return nil, err
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
	err := s.repo.Store.UpdateStore(ctx, storeID, sqlxRepo.UpdateStoreParams{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update store", err)
	}

	response := &adminStoreModel.UpdateResponse{
		ID: utils.FormatID(storeID),
	}

	return response, nil
}
