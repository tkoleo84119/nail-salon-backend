package adminProduct

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
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

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminProductModel.GetAllParsedRequest, creatorStoreIDs []int64) (*adminProductModel.GetAllResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	total, products, err := s.repo.Product.GetAllStoreProductsByFilter(ctx, storeID, sqlxRepo.GetAllStoreProductsByFilterParams{
		BrandID:             req.BrandID,
		CategoryID:          req.CategoryID,
		Name:                req.Name,
		LessThanSafetyStock: req.LessThanSafetyStock,
		IsActive:            req.IsActive,
		Limit:               &req.Limit,
		Offset:              &req.Offset,
		Sort:                &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get products", err)
	}

	items := make([]adminProductModel.GetAllProductItem, len(products))
	for i, product := range products {
		items[i] = adminProductModel.GetAllProductItem{
			ID:   utils.FormatID(product.ID),
			Name: product.Name,
			Brand: adminProductModel.GetAllProductBrandItem{
				ID:   utils.FormatID(product.BrandID),
				Name: product.BrandName,
			},
			Category: adminProductModel.GetAllProductCategoryItem{
				ID:   utils.FormatID(product.CategoryID),
				Name: product.CategoryName,
			},
			CurrentStock:    int(product.CurrentStock),
			SafetyStock:     int(utils.PgInt4ToInt32(product.SafetyStock)),
			Unit:            utils.PgTextToString(product.Unit),
			StorageLocation: utils.PgTextToString(product.StorageLocation),
			Note:            utils.PgTextToString(product.Note),
			IsActive:        utils.PgBoolToBool(product.IsActive),
			CreatedAt:       utils.PgTimestamptzToTimeString(product.CreatedAt),
			UpdatedAt:       utils.PgTimestamptzToTimeString(product.UpdatedAt),
		}
	}

	return &adminProductModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
