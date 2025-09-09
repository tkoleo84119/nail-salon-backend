package auth

import "github.com/tkoleo84119/nail-salon-backend/internal/model/common"

type LineLoginRequest struct {
	IdToken string `json:"idToken" binding:"required,noBlank,max=2000"`
}

type LineLoginResponse struct {
	NeedRegister   bool                `json:"needRegister"`
	NeedCheckTerms *bool               `json:"needCheckTerms,omitempty"`
	AccessToken    *string             `json:"accessToken,omitempty"`
	RefreshToken   *string             `json:"refreshToken,omitempty"`
	ExpiresIn      *int                `json:"expiresIn,omitempty"`
	LineProfile    *common.LineProfile `json:"lineProfile,omitempty"`
}
