package adminStockUsages

type GetAllRequest struct {
	Name    *string `form:"name" binding:"omitempty,noBlank,max=100"`
	IsInUse *bool   `form:"isInUse" binding:"omitempty"`
	Limit   *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset  *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort    *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name    *string
	IsInUse *bool
	Limit   int
	Offset  int
	Sort    []string
}

type GetAllResponse struct {
	Total int              `json:"total"`
	Items []GetAllListItem `json:"items"`
}

type GetAllListItem struct {
	ID           string            `json:"id"`
	Product      GetAllProductItem `json:"product"`
	Quantity     int32             `json:"quantity"`
	IsInUse      bool              `json:"isInUse"`
	Expiration   string            `json:"expiration"`
	UsageStarted string            `json:"usageStarted"`
	UsageEndedAt string            `json:"usageEndedAt"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

type GetAllProductItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
