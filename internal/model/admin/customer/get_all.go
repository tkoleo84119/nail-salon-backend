package adminCustomer

type GetAllRequest struct {
	Name          *string `form:"name" binding:"omitempty,max=100"`
	Phone         *string `form:"phone" binding:"omitempty,max=20"`
	Level         *string `form:"level" binding:"omitempty,oneof=NORMAL VIP VVIP"`
	IsBlacklisted *bool   `form:"isBlacklisted" binding:"omitempty,boolean"`
	MinPastDays   *int    `form:"minPastDays" binding:"omitempty,min=0,max=365"`
	Limit         *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset        *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort          *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name          *string
	Phone         *string
	Level         *string
	IsBlacklisted *bool
	MinPastDays   *int
	Limit         int
	Offset        int
	Sort          []string
}

type GetAllResponse struct {
	Total int                  `json:"total"`
	Items []GetAllCustomerItem `json:"items"`
}

type GetAllCustomerItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Birthday      string `json:"birthday"`
	City          string `json:"city"`
	Level         string `json:"level"`
	IsBlacklisted bool   `json:"isBlacklisted"`
	LastVisitAt   string `json:"lastVisitAt,omitempty"`
	UpdatedAt     string `json:"updatedAt"`
}
