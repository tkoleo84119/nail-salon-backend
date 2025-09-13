package adminCheckout

type CreateBulkRequest struct {
	PaymentMethod    string                    `json:"paymentMethod" binding:"required,oneof=cash linePay"`
	CustomerCouponID *string                   `json:"customerCouponId" binding:"omitempty"`
	Checkouts        []CreateBulkCheckoutItems `json:"checkouts" binding:"required,min=1,max=10"`
}

type CreateBulkCheckoutItems struct {
	BookingID  string                  `json:"bookingId" binding:"required"`
	PaidAmount int64                   `json:"paidAmount" binding:"required,min=0,max=1000000"`
	Details    []CreateBulkDetailItems `json:"details" binding:"required,min=1,max=10"`
}

type CreateBulkDetailItems struct {
	ID        string `json:"id" binding:"required"`
	Price     int64  `json:"price" binding:"required,min=0,max=1000000"`
	UseCoupon *bool  `json:"useCoupon" binding:"omitempty"`
}

type CreateBulkParsedRequest struct {
	PaymentMethod    string
	CustomerCouponID *int64
	Checkouts        []CreateBulkParsedCheckoutItems
}

type CreateBulkParsedCheckoutItems struct {
	BookingID  int64
	PaidAmount int64
	ApplyCount int64
	Details    []CreateBulkParsedDetailItems
}

type CreateBulkParsedDetailItems struct {
	ID        int64
	Price     int64
	UseCoupon bool
}

type CreateBulkResponse struct {
	IDs []string `json:"ids"`
}
