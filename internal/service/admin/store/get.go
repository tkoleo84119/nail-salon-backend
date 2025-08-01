package adminStore

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreService struct {
	repo *sqlxRepo.Repositories
}

func NewGetStoreService(repo *sqlxRepo.Repositories) *GetStoreService {
	return &GetStoreService{
		repo: repo,
	}
}

func (s *GetStoreService) GetStore(ctx context.Context, storeID int64) (*adminStoreModel.GetStoreResponse, error) {
	// Get store information
	store, err := s.repo.Store.GetStore(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Build response
	response := &adminStoreModel.GetStoreResponse{
		ID:        utils.FormatID(store.ID),
		Name:      store.Name,
		Address:   utils.PgTextToString(store.Address),
		Phone:     utils.PgTextToString(store.Phone),
		IsActive:  utils.PgBoolToBool(store.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(store.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(store.UpdatedAt),
	}

	return response, nil
}
