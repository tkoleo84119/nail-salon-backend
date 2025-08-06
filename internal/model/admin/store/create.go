package adminStore

type CreateRequest struct {
	Name    string  `json:"name" binding:"required,max=100"`
	Address *string `json:"address,omitempty" binding:"omitempty,max=255"`
	Phone   *string `json:"phone,omitempty" binding:"omitempty,max=20,taiwanlandline"`
}

type CreateResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"isActive"`
}
