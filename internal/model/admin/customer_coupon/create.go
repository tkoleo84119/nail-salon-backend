package adminCustomerCoupon

type CreateRequest struct {
	CustomerId string `json:"customerId" binding:"required"`
	CouponId   string `json:"couponId" binding:"required"`
	Period     string `json:"period" binding:"required,oneof=unlimited 1month 3months 6months 1year"`
}

type CreateParsedRequest struct {
	CustomerId int64
	CouponId   int64
	Period     string
}

type CreateResponse struct {
	ID string `json:"id"`
}
