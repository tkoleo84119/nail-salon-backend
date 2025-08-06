package adminStaff

type GetMeResponse struct {
	ID        string               `json:"id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Role      string               `json:"role"`
	IsActive  bool                 `json:"isActive"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
	Stylist   *GetStaffStylistInfo `json:"stylist,omitempty"`
}
