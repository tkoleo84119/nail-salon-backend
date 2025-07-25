package store

// GetStoreServicesQueryParams represents query parameters for getting store services
type GetStoreServicesQueryParams struct {
	Limit   int  `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Offset  int  `form:"offset,default=0" binding:"omitempty,min=0"`
	IsAddon *bool `form:"isAddon" binding:"omitempty"`
}

// GetStoreServicesResponse represents the response for getting store services
type GetStoreServicesResponse struct {
	Total int                        `json:"total"`
	Items []GetStoreServicesItemModel `json:"items"`
}

// GetStoreServicesItemModel represents a single service item in the store services list
type GetStoreServicesItemModel struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	DurationMinutes int     `json:"durationMinutes"`
	IsAddon         bool    `json:"isAddon"`
	Note            *string `json:"note,omitempty"`
}