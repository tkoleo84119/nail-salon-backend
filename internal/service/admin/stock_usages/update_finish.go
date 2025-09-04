package adminStockUsages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateFinish struct {
	queries *dbgen.Queries
}

func NewUpdateFinish(queries *dbgen.Queries) *UpdateFinish {
	return &UpdateFinish{
		queries: queries,
	}
}

func (s *UpdateFinish) UpdateFinish(ctx context.Context, storeID int64, stockUsageID int64, req adminStockUsagesModel.UpdateFinishParsedRequest, staffStoreIDs []int64) (*adminStockUsagesModel.UpdateFinishResponse, error) {
	if err := utils.CheckStoreAccess(storeID, staffStoreIDs); err != nil {
		return nil, err
	}

	// Check if stock usage exists
	stockUsage, err := s.queries.GetStockUsageByID(ctx, stockUsageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StockUsageNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get stock usage", err)
	}

	// Get product to check store access
	product, err := s.queries.GetProductByID(ctx, stockUsage.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get product", err)
	}
	// Check if product belongs to the specified store
	if product.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StockUsageNotFound)
	}

	// Check if stock usage is currently in use
	if !stockUsage.IsInUse.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StockUsageNotInUse)
	}

	// Update stock usage
	err = s.queries.UpdateStockUsageFinish(ctx, dbgen.UpdateStockUsageFinishParams{
		ID:           stockUsageID,
		UsageEndedAt: utils.TimeToPgDate(req.UsageEndedAt),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update stock usage", err)
	}

	return &adminStockUsagesModel.UpdateFinishResponse{
		ID: utils.FormatID(stockUsageID),
	}, nil
}
