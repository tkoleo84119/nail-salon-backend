package store

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type GetStoresService struct {
	repo *sqlxRepo.Repositories
}

func NewGetStoresService(repo *sqlxRepo.Repositories) GetStoresServiceInterface {
	return &GetStoresService{
		repo: repo,
	}
}

func (s *GetStoresService) GetStores(ctx context.Context, queryParams storeModel.GetStoresQueryParams) (*storeModel.GetStoresResponse, error) {
	// Get stores from repository with pagination
	stores, total, err := s.repo.Store.GetStores(ctx, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stores", err)
	}

	return &storeModel.GetStoresResponse{
		Total: total,
		Items: stores,
	}, nil
}