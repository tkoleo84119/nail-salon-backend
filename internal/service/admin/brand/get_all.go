package adminBrand

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBrandModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/brand"
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

func (s *GetAll) GetAll(ctx context.Context, req adminBrandModel.GetAllParsedRequest) (*adminBrandModel.GetAllResponse, error) {
	total, items, err := s.repo.Brand.GetAllBrandsByFilter(ctx, sqlxRepo.GetAllBrandsByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get brand list", err)
	}

	itemsDTO := make([]adminBrandModel.GetAllBrandItem, len(items))
	for i, item := range items {
		itemsDTO[i] = adminBrandModel.GetAllBrandItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			IsActive:  utils.PgBoolToBool(item.IsActive),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	response := &adminBrandModel.GetAllResponse{
		Total: total,
		Items: itemsDTO,
	}

	return response, nil
}
