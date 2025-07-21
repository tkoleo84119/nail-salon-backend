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

// HasUpdates checks if at least one field is provided for update
func (req *UpdateMyCustomerRequest) HasUpdates() bool {
	return req.Name != nil || req.Phone != nil || req.Birthday != nil ||
		req.City != nil || req.FavoriteShapes != nil || req.FavoriteColors != nil ||
		req.FavoriteStyles != nil || req.IsIntrovert != nil || req.CustomerNote != nil
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