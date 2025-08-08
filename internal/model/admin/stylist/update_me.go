package adminStylist

type UpdateMeRequest struct {
	Name         *string   `json:"name" binding:"omitempty,max=50"`
	GoodAtShapes *[]string `json:"goodAtShapes" binding:"omitempty,max=20"`
	GoodAtColors *[]string `json:"goodAtColors" binding:"omitempty,max=20"`
	GoodAtStyles *[]string `json:"goodAtStyles" binding:"omitempty,max=20"`
	IsIntrovert  *bool     `json:"isIntrovert" binding:"omitempty,boolean"`
}

type UpdateMeResponse struct {
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

func (r *UpdateMeRequest) HasUpdate() bool {
	return r.Name != nil || r.GoodAtShapes != nil || r.GoodAtColors != nil || r.GoodAtStyles != nil || r.IsIntrovert != nil
}
