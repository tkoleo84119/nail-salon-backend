package adminProduct

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) GetInterface {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID, productID int64, role string, creatorStoreIDs []int64) (*adminProductModel.GetResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	product, err := s.queries.GetProductWithDetailsByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get product", err)
	}

	// check product is belong to the store
	if product.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
	}

	return &adminProductModel.GetResponse{
		ID:   utils.FormatID(product.ID),
		Name: product.Name,
		Brand: adminProductModel.GetProductBrandItem{
			ID:   utils.FormatID(product.BrandID),
			Name: product.BrandName,
		},
		Category: adminProductModel.GetProductCategoryItem{
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
	}, nil
}
