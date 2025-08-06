package adminStaff

type UpdateMeRequest struct {
	Email *string `json:"email" binding:"omitempty,email"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateMeRequest) HasUpdates() bool {
	return r.Email != nil
}

// UpdateMyStaffResponse represents the response after updating current staff user's information
type UpdateMeResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
