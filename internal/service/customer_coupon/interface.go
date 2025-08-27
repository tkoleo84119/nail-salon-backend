package customerCoupon

import (
	"context"

	customerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer_coupon"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, customerID int64, req customerCouponModel.GetAllParsedRequest) (*customerCouponModel.GetAllResponse, error)
}
