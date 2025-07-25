package customer

import (
	"context"

	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
)

type UpdateMyCustomerServiceInterface interface {
	UpdateMyCustomer(ctx context.Context, customerID int64, req customerModel.UpdateMyCustomerRequest) (*customerModel.UpdateMyCustomerResponse, error)
}
