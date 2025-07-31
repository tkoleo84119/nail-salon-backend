package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreListService struct {
	repo *sqlxRepo.Repositories
}

func NewGetStoreListService(repo *sqlxRepo.Repositories) *GetStoreListService {
	return &GetStoreListService{
		repo: repo,
	}
}

func (s *GetStoreListService) GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListParsedRequest) (*adminStoreModel.GetStoreListResponse, error) {
	total, items, err := s.repo.Store.GetAllByFilter(ctx, sqlxRepo.GetAllByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store list", err)
	}

	itemsDTO := make([]adminStoreModel.StoreListItemDTO, len(items))
	for i, item := range items {
		itemsDTO[i] = adminStoreModel.StoreListItemDTO{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			Address:   utils.PgTextToString(item.Address),
			Phone:     utils.PgTextToString(item.Phone),
			IsActive:  item.IsActive.Bool,
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	response := &adminStoreModel.GetStoreListResponse{
		Total: total,
		Items: itemsDTO,
	}

	return response, nil
}
