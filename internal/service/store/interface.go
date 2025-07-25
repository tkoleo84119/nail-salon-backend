package store

import (
	"context"

	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
)

type GetStoreServicesServiceInterface interface {
	GetStoreServices(ctx context.Context, storeIDStr string, queryParams storeModel.GetStoreServicesQueryParams) (*storeModel.GetStoreServicesResponse, error)
}

type GetStoreStylistsServiceInterface interface {
	GetStoreStylists(ctx context.Context, storeIDStr string, queryParams storeModel.GetStoreStylistsQueryParams) (*storeModel.GetStoreStylistsResponse, error)
}

type GetStoresServiceInterface interface {
	GetStores(ctx context.Context, queryParams storeModel.GetStoresQueryParams) (*storeModel.GetStoresResponse, error)
}