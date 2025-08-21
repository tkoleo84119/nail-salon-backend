package adminService

// GetServiceListRequest represents the request to get service list with filtering
type GetAllRequest struct {
	Name      *string `form:"name" binding:"omitempty,max=100"`
	IsAddon   *bool   `form:"isAddon" binding:"omitempty"`
	IsActive  *bool   `form:"isActive" binding:"omitempty"`
	IsVisible *bool   `form:"isVisible" binding:"omitempty"`
	Limit     *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort      *string `form:"sort" binding:"omitempty"`
}

type GetAllParsedRequest struct {
	Name      *string
	IsAddon   *bool
	IsActive  *bool
	IsVisible *bool
	Limit     int
	Offset    int
	Sort      []string
}

// GetServiceListResponse represents the response for service list
type GetAllResponse struct {
	Total int                        `json:"total"`
	Items []GetAllServiceListItemDTO `json:"items"`
}

// ServiceListItemDTO represents a single service item in the list
type GetAllServiceListItemDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int32  `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	IsActive        bool   `json:"isActive"`
	IsVisible       bool   `json:"isVisible"`
	Note            string `json:"note"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}
