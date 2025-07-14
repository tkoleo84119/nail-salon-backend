package staff

import (
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

// LoginRequest represents the staff login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=1,max=100"`
	Password string `json:"password" binding:"required,min=1,max=100"`
}

// LoginResponse represents the successful login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

// User represents the authenticated staff user info
type User struct {
	ID        string        `json:"id"` // Snowflake ID as string
	Username  string        `json:"username"`
	Role      string        `json:"role"`
	StoreList []common.Store `json:"store_list"`
}

// LoginContext contains the context information for login
type LoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

// TokenInfo represents staff token information
type TokenInfo struct {
	StaffUserID  int64
	RefreshToken string
	Context      LoginContext
	ExpiresAt    time.Time
}
