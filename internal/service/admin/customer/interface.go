package customer

import (
	"context"

	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminCustomerModel.GetAllParsedRequest) (*adminCustomerModel.GetAllResponse, error)
}
