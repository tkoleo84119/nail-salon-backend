package adminSupplier

import (
	"context"

	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminSupplierModel.CreateRequest) (*adminSupplierModel.CreateResponse, error)
}
