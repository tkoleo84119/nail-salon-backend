package adminStockUsages

import (
	"context"

	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminStockUsagesModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminStockUsagesModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminStockUsagesModel.GetAllParsedRequest, staffStoreIDs []int64) (*adminStockUsagesModel.GetAllResponse, error)
}
