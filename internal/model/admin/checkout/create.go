package adminCheckout

type CreateRequest struct {
	PaymentMethod    string                `json:"paymentMethod" binding:"required,oneof=cash linePay"`
	CustomerCouponID *string               `json:"customerCouponId"`
	PaidAmount       int64                 `json:"paidAmount" binding:"required"`
	BookingDetails   []CreateBookingDetail `json:"bookingDetails" binding:"required,min=1,max=10"`
}

type CreateParsedRequest struct {
	PaymentMethod    string
	CustomerCouponID *int64
	PaidAmount       int64
	BookingDetails   []CreateBookingDetailParsed
}

type CreateBookingDetail struct {
	ID        string `json:"id" binding:"required"`
	Price     int64  `json:"price" binding:"required"`
	UseCoupon bool   `json:"useCoupon" binding:"required"`
}

type CreateBookingDetailParsed struct {
	ID        int64
	Price     int64
	UseCoupon bool
}

type CreateResponse struct {
	ID string `json:"id"`
}
