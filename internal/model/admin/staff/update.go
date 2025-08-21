package adminStaff

type UpdateRequest struct {
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=ADMIN MANAGER STYLIST"`
	IsActive *bool   `json:"isActive,omitempty" binding:"omitempty"`
}

type UpdateResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.Role != nil || r.IsActive != nil
}
