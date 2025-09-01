package adminBrand

import (
	"context"

	adminBrandModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/brand"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminBrandModel.CreateRequest) (*adminBrandModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminBrandModel.GetAllParsedRequest) (*adminBrandModel.GetAllResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, brandID int64, req adminBrandModel.UpdateRequest) (*adminBrandModel.UpdateResponse, error)
}
