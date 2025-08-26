package adminCustomerCoupon

import "time"

type CreateRequest struct {
	CustomerId string  `json:"customerId" binding:"required"`
	CouponId   string  `json:"couponId" binding:"required"`
	ValidFrom  string  `json:"validFrom" binding:"required"`
	ValidTo    *string `json:"validTo" binding:"omitempty"`
}

type CreateParsedRequest struct {
	CustomerId int64
	CouponId   int64
	ValidFrom  time.Time
	ValidTo    *time.Time
}

type CreateResponse struct {
	ID string `json:"id"`
}
