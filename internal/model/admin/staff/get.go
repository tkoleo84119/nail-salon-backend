package adminStaff

type GetResponse struct {
	ID        string               `json:"id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Role      string               `json:"role"`
	IsActive  bool                 `json:"isActive"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
	Stylist   *GetStaffStylistInfo `json:"stylist"`
}

type GetStaffStylistInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}
