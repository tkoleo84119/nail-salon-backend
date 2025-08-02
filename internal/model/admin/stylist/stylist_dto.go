package adminStylist

// UpdateMyStylistRequest represents the request to update a stylist
type UpdateMyStylistRequest struct {
	Name         *string   `json:"name" binding:"omitempty,max=50"`
	GoodAtShapes *[]string `json:"goodAtShapes" binding:"omitempty,max=20"`
	GoodAtColors *[]string `json:"goodAtColors" binding:"omitempty,max=20"`
	GoodAtStyles *[]string `json:"goodAtStyles" binding:"omitempty,max=20"`
	IsIntrovert  *bool     `json:"isIntrovert" binding:"omitempty,boolean"`
}

// UpdateMyStylistResponse represents the response after updating a stylist
type UpdateMyStylistResponse struct {
	ID           string   `json:"id"`
	StaffUserID  string   `json:"staffUserId"`
	Name         string   `json:"name"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

func (r *UpdateMyStylistRequest) HasUpdate() bool {
	return r.Name != nil || r.GoodAtShapes != nil || r.GoodAtColors != nil || r.GoodAtStyles != nil || r.IsIntrovert != nil
}

// -------------------------------------------------------------------------------------

// GetStylistListRequest represents the request to get stylists list
type GetStylistListRequest struct {
	Name        *string `form:"name" binding:"omitempty,max=50"`
	IsIntrovert *bool   `form:"isIntrovert"`
	Limit       *int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset      *int    `form:"offset" binding:"omitempty,min=0"`
}

// GetStylistListResponse represents the response with stylists list
type GetStylistListResponse struct {
	Total int                  `json:"total"`
	Items []GetStylistListItem `json:"items"`
}

// GetStylistListItem represents a single stylist in the list
type GetStylistListItem struct {
	ID           string   `json:"id"`
	StaffUserID  string   `json:"staffUserId"`
	Name         string   `json:"name"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
}
