package store

type GetAllRequest struct {
	Limit  *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort   *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Limit  int
	Offset int
	Sort   []string
}

type GetAllResponse struct {
	Total int                   `json:"total"`
	Items []GetAllStoreListItem `json:"items"`
}

type GetAllStoreListItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}
