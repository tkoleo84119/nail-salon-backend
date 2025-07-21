package customer

// LineRegisterRequest represents the request for LINE customer registration
type LineRegisterRequest struct {
	IdToken         string    `json:"idToken" binding:"required,min=1,max=500"`
	Name            string    `json:"name" binding:"required,min=1,max=100"`
	Phone           string    `json:"phone" binding:"required,min=1,max=20,taiwanmobile"`
	Birthday        string    `json:"birthday" binding:"required"`
	City            string    `json:"city,omitempty" binding:"omitempty,min=1,max=100"`
	FavoriteShapes  []string  `json:"favorite_shapes,omitempty" binding:"omitempty,min=1,max=20"`
	FavoriteColors  []string  `json:"favorite_colors,omitempty" binding:"omitempty,min=1,max=20"`
	FavoriteStyles  []string  `json:"favorite_styles,omitempty" binding:"omitempty,min=1,max=20"`
	IsIntrovert     *bool     `json:"is_introvert,omitempty"`
	ReferralSource  []string  `json:"referral_source,omitempty" binding:"omitempty,min=1,max=20"`
	Referrer        string    `json:"referrer,omitempty" binding:"omitempty,min=1,max=100"`
	CustomerNote    string    `json:"customer_note,omitempty" binding:"omitempty,min=1,max=1000"`
}

// LineRegisterResponse represents the response for LINE customer registration
type LineRegisterResponse struct {
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	Customer     *RegisteredCustomer `json:"customer"`
}

// RegisteredCustomer represents the customer information after registration
type RegisteredCustomer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
}