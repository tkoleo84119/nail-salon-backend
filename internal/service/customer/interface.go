package customer

import (
	"context"

	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
)

type GetMeInterface interface {
	GetMe(ctx context.Context, customerID int64) (*customerModel.GetMeResponse, error)
}

type UpdateMeInterface interface {
	UpdateMe(ctx context.Context, customerID int64, req customerModel.UpdateMeRequest) (*customerModel.UpdateMeResponse, error)
}
