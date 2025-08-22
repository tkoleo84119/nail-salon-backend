package adminStore

import (
	"context"
	"strings"

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
	// For ADMIN role, check store access permission
	if role == common.RoleAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
			return nil, err
		}
	}

	// Check if name is unique (excluding current store)
	var name string
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
		nameExists, err := s.queries.CheckStoreNameExistsExcluding(ctx, dbgen.CheckStoreNameExistsExcludingParams{
			ID:   storeID,
			Name: name,
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
		Name:     &name,
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
