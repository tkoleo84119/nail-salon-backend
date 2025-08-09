package auth

import "time"

type LineRegisterRequest struct {
	IdToken        string    `json:"idToken" binding:"required,max=500"`
	Name           string    `json:"name" binding:"required,max=100"`
	Phone          string    `json:"phone" binding:"required,taiwanmobile"`
	Birthday       string    `json:"birthday" binding:"required"`
	City           *string   `json:"city" binding:"omitempty,max=100"`
	FavoriteShapes *[]string `json:"favoriteShapes" binding:"omitempty,max=20"`
	FavoriteColors *[]string `json:"favoriteColors" binding:"omitempty,max=20"`
	FavoriteStyles *[]string `json:"favoriteStyles" binding:"omitempty,max=20"`
	IsIntrovert    *bool     `json:"isIntrovert" binding:"omitempty,boolean"`
	ReferralSource *[]string `json:"referralSource" binding:"omitempty,max=20"`
	Referrer       *string   `json:"referrer" binding:"omitempty,max=100"`
	CustomerNote   *string   `json:"customerNote" binding:"omitempty,max=255"`
}

type LoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

type LineRegisterResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}
