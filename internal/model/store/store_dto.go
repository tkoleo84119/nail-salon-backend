package store

// CreateStoreRequest represents the request to create a new store
type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=100"`
	Address string `json:"address" binding:"omitempty,max=255"`
	Phone   string `json:"phone" binding:"omitempty,taiwanlandline"`
}

// CreateStoreResponse represents the response after creating a store
type CreateStoreResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"isActive"`
}
