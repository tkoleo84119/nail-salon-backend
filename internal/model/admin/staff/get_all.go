package adminStaff

type GetAllRequest struct {
	Username *string `form:"username" binding:"omitempty,max=100"`
	Email    *string `form:"email" binding:"omitempty,max=100"`
	Role     *string `form:"role," binding:"omitempty,oneof=SUPER_ADMIN ADMIN MANAGER STYLIST"`
	IsActive *bool   `form:"isActive" binding:"omitempty,boolean"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Username *string
	Email    *string
	Role     *string
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

type GetAllResponse struct {
	Total int                   `json:"total"`
	Items []GetAllStaffListItem `json:"items"`
}

type GetAllStaffListItem struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
