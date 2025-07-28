package adminService

// CreateServiceRequest represents the request to create a new service
type CreateServiceRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Price           int64  `json:"price" binding:"required,min=0"`
	DurationMinutes int32  `json:"durationMinutes" binding:"required,min=0,max=1440"`
	IsAddon         bool   `json:"isAddon" binding:"omitempty"`
	IsVisible       bool   `json:"isVisible" binding:"omitempty"`
	Note            string `json:"note" binding:"omitempty,max=255"`
}

// CreateServiceResponse represents the response after creating a service
type CreateServiceResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int32  `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	IsVisible       bool   `json:"isVisible"`
	IsActive        bool   `json:"isActive"`
	Note            string `json:"note"`
}

// -------------------------------------------------------------------------------------

// UpdateServiceRequest represents the request to update a service
type UpdateServiceRequest struct {
	Name            *string `json:"name" binding:"omitempty,min=1,max=100"`
	Price           *int64  `json:"price" binding:"omitempty,min=0"`
	DurationMinutes *int32  `json:"durationMinutes" binding:"omitempty,min=0,max=1440"`
	IsAddon         *bool   `json:"isAddon" binding:"omitempty"`
	IsVisible       *bool   `json:"isVisible" binding:"omitempty"`
	IsActive        *bool   `json:"isActive" binding:"omitempty"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

// UpdateServiceResponse represents the response after updating a service
type UpdateServiceResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int32  `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	IsVisible       bool   `json:"isVisible"`
	IsActive        bool   `json:"isActive"`
	Note            string `json:"note"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateServiceRequest) HasUpdates() bool {
	return r.Name != nil || r.Price != nil || r.DurationMinutes != nil ||
		r.IsAddon != nil || r.IsVisible != nil || r.IsActive != nil || r.Note != nil
}

// -------------------------------------------------------------------------------------

// GetServiceListRequest represents the request to get service list with filtering
type GetServiceListRequest struct {
	Name      *string `form:"name" binding:"omitempty,max=100"`
	IsAddon   *bool   `form:"isAddon"`
	IsActive  *bool   `form:"isActive"`
	IsVisible *bool   `form:"isVisible"`
	Limit     *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    *int    `form:"offset" binding:"omitempty,min=0"`
}

// GetServiceListResponse represents the response for service list
type GetServiceListResponse struct {
	Total int                  `json:"total"`
	Items []ServiceListItemDTO `json:"items"`
}

// ServiceListItemDTO represents a single service item in the list
type ServiceListItemDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int32  `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	IsActive        bool   `json:"isActive"`
	IsVisible       bool   `json:"isVisible"`
	Note            string `json:"note"`
}
