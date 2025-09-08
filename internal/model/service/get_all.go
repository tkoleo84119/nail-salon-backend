package service

type GetAllRequest struct {
	IsAddon *bool   `form:"isAddon" binding:"omitempty"`
	Limit   *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset  *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort    *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	IsAddon *bool
	Limit   int
	Offset  int
	Sort    []string
}

type GetAllResponse struct {
	Total int          `json:"total"`
	Items []GetAllItem `json:"items"`
}

type GetAllItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int    `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	Note            string `json:"note"`
}
