package adminStore

type GetAllRequest struct {
	Name     *string `form:"name" binding:"omitempty,max=100"`
	IsActive *bool   `form:"isActive" binding:"omitempty,boolean"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name     *string
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

type GetAllResponse struct {
	Total int                   `json:"total"`
	Items []GetAllStoreListItem `json:"items"`
}

type GetAllStoreListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
