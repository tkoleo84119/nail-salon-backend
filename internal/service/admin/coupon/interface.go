package adminCoupon

import (
	"context"

	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminCouponModel.CreateRequest) (*adminCouponModel.CreateResponse, error)
}
