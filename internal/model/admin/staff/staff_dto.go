package adminStaff

// CreateStaffRequest represents the request to create a new staff member
type CreateStaffRequest struct {
	Username string   `json:"username" binding:"required,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,max=50"`
	Role     string   `json:"role" binding:"required,oneof=ADMIN MANAGER STYLIST"`
	StoreIDs []string `json:"storeIds" binding:"required,min=1,max=10"`
}

type CreateStaffParsedRequest struct {
	Username string
	Email    string
	Password string
	Role     string
	StoreIDs []int64
}

// CreateStaffResponse represents the response after creating a staff member
type CreateStaffResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// UpdateStaffRequest represents the request to update staff information
type UpdateStaffRequest struct {
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=ADMIN MANAGER STYLIST"`
	IsActive *bool   `json:"isActive,omitempty" binding:"omitempty,boolean"`
}

// UpdateStaffResponse represents the response after updating staff information
type UpdateStaffResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (r UpdateStaffRequest) HasUpdates() bool {
	return r.Role != nil || r.IsActive != nil
}

// -------------------------------------------------------------------------------------

// UpdateMyStaffRequest represents the request to update current staff user's information
type UpdateMyStaffRequest struct {
	Email *string `json:"email" binding:"omitempty,email"`
}

// HasUpdates checks if the request has any fields to update
func (r UpdateMyStaffRequest) HasUpdates() bool {
	return r.Email != nil
}

// UpdateMyStaffResponse represents the response after updating current staff user's information
type UpdateMyStaffResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// GetStaffListRequest represents the request to get staff list with filtering
type GetStaffListRequest struct {
	Username *string `form:"username" binding:"omitempty,max=100"`
	Email    *string `form:"email" binding:"omitempty,max=100"`
	Role     *string `form:"role," binding:"omitempty,oneof=SUPER_ADMIN ADMIN MANAGER STYLIST"`
	IsActive *bool   `form:"isActive" binding:"omitempty,boolean"`
	Limit    *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset   *int    `form:"offset" binding:"omitempty,min=0,max=1000000"`
	Sort     *string `form:"sort" binding:"omitempty"`
}

type GetStaffListParsedRequest struct {
	Username *string
	Email    *string
	Role     *string
	IsActive *bool
	Limit    int
	Offset   int
	Sort     []string
}

// GetStaffListResponse represents the response for staff list
type GetStaffListResponse struct {
	Total int                `json:"total"`
	Items []StaffListItemDTO `json:"items"`
}

// StaffListItemDTO represents a single staff item in the list
type StaffListItemDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// -------------------------------------------------------------------------------------

// GetMyStaffResponse represents the response for current staff member's information
type GetMyStaffResponse struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Role      string            `json:"role"`
	IsActive  bool              `json:"isActive"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
	Stylist   *StaffStylistInfo `json:"stylist,omitempty"`
}

// -------------------------------------------------------------------------------------

// GetStaffResponse represents the response for a specific staff member with optional stylist information
type GetStaffResponse struct {
	ID        string            `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Role      string            `json:"role"`
	IsActive  bool              `json:"isActive"`
	CreatedAt string            `json:"createdAt"`
	UpdatedAt string            `json:"updatedAt"`
	Stylist   *StaffStylistInfo `json:"stylist"`
}

// StaffStylistInfo represents stylist information within staff response
type StaffStylistInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}
