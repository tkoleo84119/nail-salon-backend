package adminStoreAccess

import (
	"context"

	adminStoreAccessModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store_access"
)

type CreateInterface interface {
	Create(ctx context.Context, staffID int64, storeID int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStoreAccessModel.CreateResponse, bool, error)
}

type GetInterface interface {
	Get(ctx context.Context, staffID int64) (*adminStoreAccessModel.GetResponse, error)
}

type DeleteBulkInterface interface {
	DeleteBulk(ctx context.Context, targetID int64, storeIDs []int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStoreAccessModel.DeleteBulkResponse, error)
}
