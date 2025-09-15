package adminStockUsages

import (
	"context"

	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminStockUsagesModel.CreateParsedRequest, role string, creatorStoreIDs []int64) (*adminStockUsagesModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminStockUsagesModel.GetAllParsedRequest, role string, staffStoreIDs []int64) (*adminStockUsagesModel.GetAllResponse, error)
}

type UpdateFinishInterface interface {
	UpdateFinish(ctx context.Context, storeID int64, stockUsageID int64, req adminStockUsagesModel.UpdateFinishParsedRequest, role string, staffStoreIDs []int64) (*adminStockUsagesModel.UpdateFinishResponse, error)
}
