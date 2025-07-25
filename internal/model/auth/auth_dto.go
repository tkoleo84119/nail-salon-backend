package auth

import (
	"time"
)

// CustomerLineLoginRequest represents the LINE login request payload
type CustomerLineLoginRequest struct {
	IdToken string `json:"idToken" binding:"required,min=1,max=500"`
}

// CustomerLineLoginResponse represents the LINE login response when customer is already registered
type CustomerLineLoginResponse struct {
	NeedRegister bool                 `json:"needRegister"`
	AccessToken  *string              `json:"accessToken,omitempty"`
	RefreshToken *string              `json:"refreshToken,omitempty"`
	Customer     *Customer            `json:"customer,omitempty"`
	LineProfile  *CustomerLineProfile `json:"lineProfile,omitempty"`
}

// Customer represents the customer basic information
type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// CustomerLineProfile represents the LINE profile information for registration
type CustomerLineProfile struct {
	ProviderUid string  `json:"providerUid"`
	Name        string  `json:"name"`
	Email       *string `json:"email,omitempty"`
}

// CustomerLoginContext contains the context information for customer login
type CustomerLoginContext struct {
	UserAgent string
	IPAddress string
	Timestamp time.Time
}

// CustomerTokenInfo represents customer token information
type CustomerTokenInfo struct {
	CustomerID   int64
	RefreshToken string
	Context      CustomerLoginContext
	ExpiresAt    time.Time
}

// CustomerProfile represents customer profile information from provider
type CustomerProfile struct {
	ProviderUid string
	Name        string
	Email       *string
	PictureURL  *string
}

// -------------------------------------------------------------------------------------

// CustomerLineRegisterRequest represents the request for LINE customer registration
type CustomerLineRegisterRequest struct {
	IdToken        string   `json:"idToken" binding:"required,min=1,max=500"`
	Name           string   `json:"name" binding:"required,min=1,max=100"`
	Phone          string   `json:"phone" binding:"required,min=1,max=20,taiwanmobile"`
	Birthday       string   `json:"birthday" binding:"required"`
	City           string   `json:"city,omitempty" binding:"omitempty,min=1,max=100"`
	FavoriteShapes []string `json:"favorite_shapes,omitempty" binding:"omitempty,min=1,max=20"`
	FavoriteColors []string `json:"favorite_colors,omitempty" binding:"omitempty,min=1,max=20"`
	FavoriteStyles []string `json:"favorite_styles,omitempty" binding:"omitempty,min=1,max=20"`
	IsIntrovert    *bool    `json:"is_introvert,omitempty"`
	ReferralSource []string `json:"referral_source,omitempty" binding:"omitempty,min=1,max=20"`
	Referrer       string   `json:"referrer,omitempty" binding:"omitempty,min=1,max=100"`
	CustomerNote   string   `json:"customer_note,omitempty" binding:"omitempty,min=1,max=1000"`
}

// CustomerLineRegisterResponse represents the response for LINE customer registration
type CustomerLineRegisterResponse struct {
	AccessToken  string                      `json:"accessToken"`
	RefreshToken string                      `json:"refreshToken"`
	Customer     *CustomerRegisteredCustomer `json:"customer"`
}

// CustomerRegisteredCustomer represents the customer information after registration
type CustomerRegisteredCustomer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
}
