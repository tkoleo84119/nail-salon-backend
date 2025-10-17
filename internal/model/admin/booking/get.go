package adminBooking

type GetResponse struct {
	ID             string                 `json:"id"`
	Customer       GetCustomer            `json:"customer"`
	Stylist        GetStylist             `json:"stylist"`
	TimeSlot       GetTimeSlot            `json:"timeSlot"`
	ActualDuration *int32                 `json:"actualDuration,omitempty"`
	Status         string                 `json:"status"`
	IsChatEnabled  bool                   `json:"isChatEnabled"`
	Note           string                 `json:"note"`
	StoreNote      string                 `json:"storeNote"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
	BookingDetails []GetBookingDetailItem `json:"bookingDetails"`
	Checkout       *GetCheckout           `json:"checkout"`
}

type GetCustomer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetStylist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetTimeSlot struct {
	ID        string `json:"id"`
	WorkDate  string `json:"workDate"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type GetBookingDetailItem struct {
	ID       string     `json:"id"`
	Service  GetService `json:"service"`
	RawPrice float64    `json:"rawPrice"`
	Price    float64    `json:"price"`
}

type GetService struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAddon bool   `json:"isAddon"`
}

type GetCheckout struct {
	ID            string     `json:"id"`
	PaymentMethod string     `json:"paymentMethod"`
	TotalAmount   int64      `json:"totalAmount"`
	FinalAmount   int64      `json:"finalAmount"`
	PaidAmount    int64      `json:"paidAmount"`
	CheckoutUser  string     `json:"checkoutUser"`
	Coupon        *GetCoupon `json:"coupon"`
}

type GetCoupon struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
