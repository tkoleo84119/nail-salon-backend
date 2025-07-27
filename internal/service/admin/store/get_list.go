package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

type GetStoreListService struct {
	repository sqlx.StoreRepositoryInterface
}

func NewGetStoreListService(repository sqlx.StoreRepositoryInterface) *GetStoreListService {
	return &GetStoreListService{
		repository: repository,
	}
}

func (s *GetStoreListService) GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error) {
	response, err := s.repository.GetStoreList(ctx, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get store list", err)
	}

	return response, nil
}