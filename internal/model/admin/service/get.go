package adminService

type GetResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DurationMinutes int32  `json:"durationMinutes"`
	Price           int64  `json:"price"`
	IsAddon         bool   `json:"isAddon"`
	IsActive        bool   `json:"isActive"`
	IsVisible       bool   `json:"isVisible"`
	Note            string `json:"note"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}
