package adminStore

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

// -------------------------------------------------------------------------------------

// UpdateStoreRequest represents the request to update a store
type UpdateStoreRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1,max=100"`
	Address  *string `json:"address" binding:"omitempty,max=255"`
	Phone    *string `json:"phone" binding:"omitempty,taiwanlandline"`
	IsActive *bool   `json:"isActive" binding:"omitempty"`
}

// UpdateStoreResponse represents the response after updating a store
type UpdateStoreResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"isActive"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateStoreRequest) HasUpdates() bool {
	return r.Name != nil || r.Address != nil || r.Phone != nil || r.IsActive != nil
}
