package customer

import (
	"time"
)

// LineLoginRequest represents the LINE login request payload
type LineLoginRequest struct {
	IdToken string `json:"idToken" binding:"required,min=1,max=500"`
}

// LineLoginResponse represents the LINE login response when customer is already registered
type LineLoginResponse struct {
	NeedRegister bool         `json:"needRegister"`
	AccessToken  *string      `json:"accessToken,omitempty"`
	RefreshToken *string      `json:"refreshToken,omitempty"`
	Customer     *Customer    `json:"customer,omitempty"`
	LineProfile  *LineProfile `json:"lineProfile,omitempty"`
}

// Customer represents the customer basic information
type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// LineProfile represents the LINE profile information for registration
type LineProfile struct {
	ProviderUid string  `json:"providerUid"`
	Name        string  `json:"name"`
	Email       *string `json:"email,omitempty"`
}

// LoginContext contains the context information for customer login
type LoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

// TokenInfo represents customer token information
type TokenInfo struct {
	CustomerID   int64
	RefreshToken string
	Context      LoginContext
	ExpiresAt    time.Time
}

// CustomerAuth represents customer authentication record
type CustomerAuth struct {
	ID          int64
	CustomerID  int64
	Provider    string
	ProviderUid string
	OtherInfo   map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CustomerProfile represents customer profile information from provider
type CustomerProfile struct {
	ProviderUid string
	Name        string
	Email       *string
	PictureURL  *string
}