package adminService

type CreateRequest struct {
	Name            string  `json:"name" binding:"required,max=100"`
	Price           *int64  `json:"price" binding:"required,min=0,max=1000000"`
	DurationMinutes *int32  `json:"durationMinutes" binding:"required,min=0,max=1440"`
	IsAddon         bool    `json:"isAddon" binding:"omitempty,boolean"`
	IsVisible       bool    `json:"isVisible" binding:"omitempty,boolean"`
	Note            *string `json:"note,omitempty" binding:"omitempty,max=255"`
}

type CreateResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	DurationMinutes int32  `json:"durationMinutes"`
	IsAddon         bool   `json:"isAddon"`
	IsVisible       bool   `json:"isVisible"`
	IsActive        bool   `json:"isActive"`
	Note            string `json:"note"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}
