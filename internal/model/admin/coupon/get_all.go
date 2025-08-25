package adminCoupon

type GetAllRequest struct {
	Name     *string `form:"name" binding:"omitempty,noBlank,max=100"`
	Code     *string `form:"code" binding:"omitempty,noBlank,max=100"`
	IsActive *bool   `form:"isActive" binding:"omitempty"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name     *string
	Code     *string
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

type GetAllResponse struct {
	Total int                   `json:"total"`
	Items []GetAllCouponItemDTO `json:"items"`
}

type GetAllCouponItemDTO struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	DisplayName    string  `json:"displayName"`
	Code           string  `json:"code"`
	DiscountRate   float64 `json:"discountRate"`
	DiscountAmount int64   `json:"discountAmount"`
	IsActive       bool    `json:"isActive"`
	Note           string  `json:"note"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}
