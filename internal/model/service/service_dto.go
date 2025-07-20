package service

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
