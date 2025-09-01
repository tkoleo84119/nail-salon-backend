package adminProductCategory

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductCategoryModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product_category"
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

func (s *GetAll) GetAll(ctx context.Context, req adminProductCategoryModel.GetAllParsedRequest) (*adminProductCategoryModel.GetAllResponse, error) {
	// 查詢資料
	total, items, err := s.repo.ProductCategory.GetAllProductCategoriesByFilter(ctx, sqlxRepo.GetAllProductCategoriesByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get product categories", err)
	}

	responseItems := make([]adminProductCategoryModel.GetAllResponseItem, len(items))
	for i, item := range items {
		responseItems[i] = adminProductCategoryModel.GetAllResponseItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			IsActive:  utils.PgBoolToBool(item.IsActive),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	return &adminProductCategoryModel.GetAllResponse{
		Total: total,
		Items: responseItems,
	}, nil
}
