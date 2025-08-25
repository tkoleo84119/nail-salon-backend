package adminAuth

import (
	"time"
)

// StaffLoginRequest represents the staff login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,noBlank,max=100"`
	Password string `json:"password" binding:"required,noBlank,max=100"`
}

// StaffLoginResponse represents the successful login response
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// StaffLoginContext contains the context information for login
type LoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

// StaffTokenInfo represents staff token information
type LoginTokenInfo struct {
	StaffUserID  int64
	RefreshToken string
	Context      LoginContext
	ExpiresAt    time.Time
}
