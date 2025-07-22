package customer

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
)

type UpdateMyCustomerServiceInterface interface {
	UpdateMyCustomer(ctx context.Context, customerID int64, req customer.UpdateMyCustomerRequest) (*customer.UpdateMyCustomerResponse, error)
}
