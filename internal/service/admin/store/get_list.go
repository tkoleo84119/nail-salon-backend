package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type GetStoreListService struct {
	repo *sqlx.Repositories
}

func NewGetStoreListService(repo *sqlx.Repositories) *GetStoreListService {
	return &GetStoreListService{
		repo: repo,
	}
}

func (s *GetStoreListService) GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error) {
	response, err := s.repo.Store.GetStoreList(ctx, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get store list", err)
	}

	return response, nil
}