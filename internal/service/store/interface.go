package store

import (
	"context"

	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, queryParams storeModel.GetAllParsedRequest) (*storeModel.GetAllResponse, error)
}

type GetStoreServicesServiceInterface interface {
	GetStoreServices(ctx context.Context, storeIDStr string, queryParams storeModel.GetStoreServicesQueryParams) (*storeModel.GetStoreServicesResponse, error)
}
