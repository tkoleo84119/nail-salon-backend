package adminProduct

import (
	"context"

	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminProductModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminProductModel.CreateResponse, error)
}
