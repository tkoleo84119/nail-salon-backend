package staff

// UpdateStaffMeRequest represents the request to update current staff user's information
type UpdateStaffMeRequest struct {
	Email *string `json:"email,omitempty" binding:"omitempty,email"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateStaffMeRequest) HasUpdates() bool {
	return r.Email != nil
}

// UpdateStaffMeResponse represents the response after updating current staff user's information
type UpdateStaffMeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}