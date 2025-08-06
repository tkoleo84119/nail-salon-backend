package adminStore

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, req adminStoreModel.GetAllParsedRequest) (*adminStoreModel.GetAllResponse, error) {
	total, items, err := s.repo.Store.GetAllStoreByFilter(ctx, sqlxRepo.GetAllStoreByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store list", err)
	}

	itemsDTO := make([]adminStoreModel.GetAllStoreListItem, len(items))
	for i, item := range items {
		itemsDTO[i] = adminStoreModel.GetAllStoreListItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			Address:   utils.PgTextToString(item.Address),
			Phone:     utils.PgTextToString(item.Phone),
			IsActive:  item.IsActive.Bool,
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	response := &adminStoreModel.GetAllResponse{
		Total: total,
		Items: itemsDTO,
	}

	return response, nil
}
