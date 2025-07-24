package customer

// UpdateMyCustomerRequest represents the request for updating customer's own profile
type UpdateMyCustomerRequest struct {
	Name           *string   `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Phone          *string   `json:"phone,omitempty" binding:"omitempty,min=1,max=20,taiwanmobile"`
	Birthday       *string   `json:"birthday,omitempty" binding:"omitempty"`
	City           *string   `json:"city,omitempty" binding:"omitempty,min=1,max=100"`
	FavoriteShapes *[]string `json:"favoriteShapes,omitempty" binding:"omitempty,max=20"`
	FavoriteColors *[]string `json:"favoriteColors,omitempty" binding:"omitempty,max=20"`
	FavoriteStyles *[]string `json:"favoriteStyles,omitempty" binding:"omitempty,max=20"`
	IsIntrovert    *bool     `json:"isIntrovert,omitempty"`
	CustomerNote   *string   `json:"customerNote,omitempty" binding:"omitempty,min=1,max=1000"`
}

// UpdateMyCustomerResponse represents the response for updating customer's own profile
type UpdateMyCustomerResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Birthday       *string   `json:"birthday,omitempty"`
	City           *string   `json:"city,omitempty"`
	FavoriteShapes *[]string `json:"favoriteShapes,omitempty"`
	FavoriteColors *[]string `json:"favoriteColors,omitempty"`
	FavoriteStyles *[]string `json:"favoriteStyles,omitempty"`
	IsIntrovert    *bool     `json:"isIntrovert,omitempty"`
	ReferralSource *[]string `json:"referralSource,omitempty"`
	Referrer       *string   `json:"referrer,omitempty"`
	CustomerNote   *string   `json:"customerNote,omitempty"`
}

// HasUpdates checks if at least one field is provided for update
func (req *UpdateMyCustomerRequest) HasUpdates() bool {
	return req.Name != nil || req.Phone != nil || req.Birthday != nil ||
		req.City != nil || req.FavoriteShapes != nil || req.FavoriteColors != nil ||
		req.FavoriteStyles != nil || req.IsIntrovert != nil || req.CustomerNote != nil
}
