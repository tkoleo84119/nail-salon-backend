package adminCustomer

import (
	"context"

	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminCustomerModel.GetAllParsedRequest) (*adminCustomerModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, customerID int64) (*adminCustomerModel.GetResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, customerID int64, req adminCustomerModel.UpdateRequest) (*adminCustomerModel.UpdateResponse, error)
}
