package adminAuth

import (
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

// StaffLoginRequest represents the staff login request payload
type StaffLoginRequest struct {
	Username string `json:"username" binding:"required,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}

// StaffLoginResponse represents the successful login response
type StaffLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	User         User   `json:"user"`
}

// User represents the authenticated staff user info
type User struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	Role      string         `json:"role"`
	StoreList []common.Store `json:"storeList"`
}

// StaffLoginContext contains the context information for login
type StaffLoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

// StaffTokenInfo represents staff token information
type StaffTokenInfo struct {
	StaffUserID  int64
	RefreshToken string
	Context      StaffLoginContext
	ExpiresAt    time.Time
}

// -------------------------------------------------------------------------------------

// StaffRefreshTokenRequest represents the request for refreshing staff access token
type StaffRefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required,max=500"`
}

// StaffRefreshTokenResponse represents the response for staff refresh token
type StaffRefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}
