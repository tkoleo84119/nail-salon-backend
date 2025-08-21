package customer

type UpdateMeRequest struct {
	Name           *string   `json:"name" binding:"omitempty,max=100"`
	Phone          *string   `json:"phone" binding:"omitempty,taiwanmobile"`
	Birthday       *string   `json:"birthday" binding:"omitempty"`
	City           *string   `json:"city" binding:"omitempty,max=100"`
	Email          *string   `json:"email" binding:"omitempty,email"`
	FavoriteShapes *[]string `json:"favoriteShapes" binding:"omitempty,max=20"`
	FavoriteColors *[]string `json:"favoriteColors" binding:"omitempty,max=20"`
	FavoriteStyles *[]string `json:"favoriteStyles" binding:"omitempty,max=20"`
	IsIntrovert    *bool     `json:"isIntrovert" binding:"omitempty"`
	CustomerNote   *string   `json:"customerNote" binding:"omitempty,max=255"`
}

type UpdateMeResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Phone          string   `json:"phone"`
	Birthday       string   `json:"birthday"`
	Email          string   `json:"email"`
	City           string   `json:"city"`
	FavoriteShapes []string `json:"favoriteShapes"`
	FavoriteColors []string `json:"favoriteColors"`
	FavoriteStyles []string `json:"favoriteStyles"`
	IsIntrovert    bool     `json:"isIntrovert"`
	CustomerNote   string   `json:"customerNote"`
}

// HasUpdates checks if at least one field is provided for update
func (req *UpdateMeRequest) HasUpdates() bool {
	return req.Name != nil || req.Phone != nil || req.Birthday != nil || req.City != nil || req.Email != nil ||
		req.FavoriteShapes != nil || req.FavoriteColors != nil || req.FavoriteStyles != nil ||
		req.IsIntrovert != nil || req.CustomerNote != nil
}
