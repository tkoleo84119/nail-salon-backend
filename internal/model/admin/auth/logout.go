package adminAuth

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,noBlank,max=500"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}
