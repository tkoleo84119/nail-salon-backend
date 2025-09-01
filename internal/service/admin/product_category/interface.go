package adminProductCategory

import (
	"context"

	adminProductCategoryModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product_category"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminProductCategoryModel.CreateRequest) (*adminProductCategoryModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminProductCategoryModel.GetAllParsedRequest) (*adminProductCategoryModel.GetAllResponse, error)
}
