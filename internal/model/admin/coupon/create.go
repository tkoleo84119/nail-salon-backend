package adminCoupon

type CreateRequest struct {
	Name           string   `json:"name" binding:"required,noBlank,max=100"`
	DisplayName    string   `json:"displayName" binding:"required,noBlank,max=100"`
	Code           string   `json:"code" binding:"required,noBlank,max=100"`
	DiscountRate   *float64 `json:"discountRate" binding:"omitempty,min=0.1,max=0.99"`
	DiscountAmount *int64   `json:"discountAmount" binding:"omitempty,min=1,max=1000000"`
	Note           *string  `json:"note" binding:"omitempty,max=255"`
}

type CreateResponse struct {
	ID string `json:"id"`
}
