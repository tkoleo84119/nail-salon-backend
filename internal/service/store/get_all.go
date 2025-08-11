package store

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, queryParams storeModel.GetAllParsedRequest) (*storeModel.GetAllResponse, error) {
	activeCondition := true
	total, stores, err := s.repo.Store.GetAllStoreByFilter(ctx, sqlxRepo.GetAllStoreByFilterParams{
		IsActive: &activeCondition,
		Limit:    &queryParams.Limit,
		Offset:   &queryParams.Offset,
		Sort:     &queryParams.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stores", err)
	}

	items := make([]storeModel.GetAllStoreListItem, len(stores))
	for i, store := range stores {
		items[i] = storeModel.GetAllStoreListItem{
			ID:      utils.FormatID(store.ID),
			Name:    store.Name,
			Address: utils.PgTextToString(store.Address),
			Phone:   utils.PgTextToString(store.Phone),
		}
	}

	return &storeModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
