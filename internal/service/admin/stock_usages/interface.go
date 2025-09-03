package adminStockUsages

import (
	"context"

	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminStockUsagesModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminStockUsagesModel.CreateResponse, error)
}
