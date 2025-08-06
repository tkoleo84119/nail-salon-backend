package adminStore

import (
	"context"

	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminStoreModel.CreateRequest, staffId int64, role string) (*adminStoreModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminStoreModel.GetAllParsedRequest) (*adminStoreModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, storeID int64) (*adminStoreModel.GetResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, storeID int64, req adminStoreModel.UpdateRequest, role string, storeIDList []int64) (*adminStoreModel.UpdateResponse, error)
}
