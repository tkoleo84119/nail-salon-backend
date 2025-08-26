package adminCustomerCoupon

type GetAllRequest struct {
	CustomerId *string `form:"customerId" binding:"omitempty"`
	CouponId   *string `form:"couponId" binding:"omitempty"`
	IsUsed     *bool   `form:"isUsed" binding:"omitempty"`
	Limit      *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort       *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	CustomerId *int64
	CouponId   *int64
	IsUsed     *bool
	Limit      int
	Offset     int
	Sort       []string
}

type GetAllResponse struct {
	Total int                        `json:"total"`
	Items []GetAllCustomerCouponItem `json:"items"`
}

type GetAllCustomerCouponItem struct {
	ID        string                `json:"id"`
	Customer  GetAllItemCustomerDTO `json:"customer"`
	Coupon    GetAllItemCouponDTO   `json:"coupon"`
	ValidFrom string                `json:"validFrom"`
	ValidTo   string                `json:"validTo"`
	IsUsed    bool                  `json:"isUsed"`
	UsedAt    string                `json:"usedAt"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
}

type GetAllItemCustomerDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LineName string `json:"lineName"`
	Phone    string `json:"phone"`
}

type GetAllItemCouponDTO struct {
	ID             string  `json:"id"`
	DisplayName    string  `json:"displayName"`
	Code           string  `json:"code"`
	DiscountRate   float64 `json:"discountRate"`
	DiscountAmount int64   `json:"discountAmount"`
	IsActive       bool    `json:"isActive"`
}
