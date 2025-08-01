package adminStore

// CreateStoreRequest represents the request to create a new store
type CreateStoreRequest struct {
	Name    string  `json:"name" binding:"required,max=100"`
	Address *string `json:"address,omitempty" binding:"omitempty,max=255"`
	Phone   *string `json:"phone,omitempty" binding:"omitempty,max=20,taiwanlandline"`
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

// -------------------------------------------------------------------------------------

// GetStoreListRequest represents the request to get store list with filtering
type GetStoreListRequest struct {
	Name     *string `form:"name" binding:"omitempty,max=100"`
	IsActive *bool   `form:"isActive" binding:"omitempty,boolean"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetStoreListParsedRequest struct {
	Name     *string
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

// GetStoreListResponse represents the response for store list
type GetStoreListResponse struct {
	Total int                `json:"total"`
	Items []StoreListItemDTO `json:"items"`
}

// StoreListItemDTO represents a single store item in the list
type StoreListItemDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// GetStoreResponse represents the response for a specific store
type GetStoreResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"isActive"`
}
