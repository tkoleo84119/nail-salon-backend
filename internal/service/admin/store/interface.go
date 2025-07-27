package adminStore

import (
	"context"

	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateStoreServiceInterface interface {
	CreateStore(ctx context.Context, req adminStoreModel.CreateStoreRequest, staffContext common.StaffContext) (*adminStoreModel.CreateStoreResponse, error)
}

type UpdateStoreServiceInterface interface {
	UpdateStore(ctx context.Context, storeID string, req adminStoreModel.UpdateStoreRequest, staffContext common.StaffContext) (*adminStoreModel.UpdateStoreResponse, error)
}

type GetStoreListServiceInterface interface {
	GetStoreList(ctx context.Context, req adminStoreModel.GetStoreListRequest) (*adminStoreModel.GetStoreListResponse, error)
}

type GetStoreServiceInterface interface {
	GetStore(ctx context.Context, storeID string) (*adminStoreModel.GetStoreResponse, error)
}
