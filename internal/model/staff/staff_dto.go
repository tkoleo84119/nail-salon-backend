package staff

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

// CreateStaffRequest represents the request to create a new staff member
type CreateStaffRequest struct {
	Username string   `json:"username" binding:"required,min=1,max=30"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=1,max=50"`
	Role     string   `json:"role" binding:"required,oneof=ADMIN MANAGER STYLIST"`
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=100"`
}

// CreateStaffResponse represents the response after creating a staff member
type CreateStaffResponse struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	IsActive  bool           `json:"isActive"`
	StoreList []common.Store `json:"storeList"`
}

// -------------------------------------------------------------------------------------

// UpdateStaffRequest represents the request to update staff information
type UpdateStaffRequest struct {
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=ADMIN MANAGER STYLIST"`
	IsActive *bool   `json:"isActive,omitempty" binding:"omitempty"`
}

// UpdateStaffResponse represents the response after updating staff information
type UpdateStaffResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

func (r UpdateStaffRequest) HasUpdates() bool {
	return r.Role != nil || r.IsActive != nil
}

// -------------------------------------------------------------------------------------

// UpdateMyStaffRequest represents the request to update current staff user's information
type UpdateMyStaffRequest struct {
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateMyStaffRequest) HasUpdates() bool {
	return r.Email != nil
}

// UpdateMyStaffResponse represents the response after updating current staff user's information
type UpdateMyStaffResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}
