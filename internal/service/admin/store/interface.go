package adminStore

import (
	"context"

	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
)

type CreateStoreServiceInterface interface {
	CreateStore(ctx context.Context, req adminStoreModel.CreateStoreRequest, staffId int64, role string) (*adminStoreModel.CreateStoreResponse, error)
}

type UpdateStoreServiceInterface interface {
	UpdateStore(ctx context.Context, storeID int64, req adminStoreModel.UpdateStoreRequest, role string, storeIDList []int64) (*adminStoreModel.UpdateStoreResponse, error)
}

type GetStoreListServiceInterface interface {
	GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListParsedRequest) (*adminStoreModel.GetStoreListResponse, error)
}

type GetStoreServiceInterface interface {
	GetStore(ctx context.Context, storeID int64) (*adminStoreModel.GetStoreResponse, error)
}
