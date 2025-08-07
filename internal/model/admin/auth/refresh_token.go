package adminAuth

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=500"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}
