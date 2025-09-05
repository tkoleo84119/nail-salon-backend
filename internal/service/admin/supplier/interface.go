package adminSupplier

import (
	"context"

	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminSupplierModel.CreateRequest) (*adminSupplierModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminSupplierModel.GetAllParsedRequest) (*adminSupplierModel.GetAllResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, id int64, req adminSupplierModel.UpdateRequest) (*adminSupplierModel.UpdateResponse, error)
}
