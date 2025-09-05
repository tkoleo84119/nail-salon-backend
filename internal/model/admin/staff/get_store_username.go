package adminStaff

type GetStoreUsernameRequest struct {
	IsActive *bool   `form:"isActive" binding:"omitempty"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetStoreUsernameParsedRequest struct {
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

type GetStoreUsernameResponse struct {
	Total int                        `json:"total"`
	Items []GetStoreUsernameListItem `json:"items"`
}

type GetStoreUsernameListItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"isActive"`
}
