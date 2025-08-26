package adminCustomerCoupon

import (
	"context"

	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, req adminCustomerCouponModel.GetAllParsedRequest) (*adminCustomerCouponModel.GetAllResponse, error)
}

type CreateInterface interface {
	Create(ctx context.Context, req adminCustomerCouponModel.CreateParsedRequest) (*adminCustomerCouponModel.CreateResponse, error)
}
