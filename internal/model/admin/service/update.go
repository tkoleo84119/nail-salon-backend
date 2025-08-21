package adminService

type UpdateRequest struct {
	Name            *string `json:"name" binding:"omitempty,min=1,max=100"`
	Price           *int64  `json:"price" binding:"omitempty,min=0,max=1000000"`
	DurationMinutes *int32  `json:"durationMinutes" binding:"omitempty,min=0,max=1440"`
	IsAddon         *bool   `json:"isAddon" binding:"omitempty"`
	IsVisible       *bool   `json:"isVisible" binding:"omitempty"`
	IsActive        *bool   `json:"isActive" binding:"omitempty"`
	Note            *string `json:"note" binding:"omitempty,max=255"`
}

type UpdateResponse struct {
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

func (r UpdateRequest) HasUpdates() bool {
	return r.Name != nil || r.Price != nil || r.DurationMinutes != nil ||
		r.IsAddon != nil || r.IsVisible != nil || r.IsActive != nil || r.Note != nil
}
