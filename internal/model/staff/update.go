package staff

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