package stylist

// UpdateStylistRequest represents the request to update a stylist
type UpdateStylistRequest struct {
	StylistName  *string   `json:"stylistName,omitempty" binding:"omitempty,min=1,max=50"`
	GoodAtShapes *[]string `json:"goodAtShapes,omitempty" binding:"omitempty,max=100"`
	GoodAtColors *[]string `json:"goodAtColors,omitempty" binding:"omitempty,max=100"`
	GoodAtStyles *[]string `json:"goodAtStyles,omitempty" binding:"omitempty,max=100"`
	IsIntrovert  *bool     `json:"isIntrovert,omitempty" binding:"omitempty"`
}

// UpdateStylistResponse represents the response after updating a stylist
type UpdateStylistResponse struct {
	ID           string   `json:"id"`
	StaffUserID  string   `json:"staffUserId"`
	StylistName  string   `json:"stylistName"`
	GoodAtShapes []string `json:"goodAtShapes"`
	GoodAtColors []string `json:"goodAtColors"`
	GoodAtStyles []string `json:"goodAtStyles"`
	IsIntrovert  bool     `json:"isIntrovert"`
}

// HasUpdates checks if the request has any fields to update
func (r *UpdateStylistRequest) HasUpdates() bool {
	return r.StylistName != nil || r.GoodAtShapes != nil || r.GoodAtColors != nil || r.GoodAtStyles != nil || r.IsIntrovert != nil
}