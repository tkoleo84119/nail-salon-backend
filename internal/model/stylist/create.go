package stylist

// CreateStylistRequest represents the request to create a new stylist
type CreateStylistRequest struct {
	StylistName  string   `json:"stylistName" binding:"required,min=1,max=50"`
	GoodAtShapes []string `json:"goodAtShapes" binding:"omitempty,max=100"`
	GoodAtColors []string `json:"goodAtColors" binding:"omitempty,max=100"`
	GoodAtStyles []string `json:"goodAtStyles" binding:"omitempty,max=100"`
	IsIntrovert  *bool    `json:"isIntrovert" binding:"omitempty"`
}

// CreateStylistResponse represents the response after creating a stylist
type CreateStylistResponse struct {
	ID           string   `json:"id"`
	StaffUserID  string   `json:"staffUserId"`
	StylistName  string   `json:"stylistName"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
}
