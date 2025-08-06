package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
)

type LoginInterface interface {
	Login(ctx context.Context, req adminAuthModel.LoginRequest, loginCtx adminAuthModel.LoginContext) (*adminAuthModel.LoginResponse, error)
}

type LogoutInterface interface {
	Logout(ctx context.Context, req adminAuthModel.LogoutRequest) (*adminAuthModel.LogoutResponse, error)
}

type RefreshTokenInterface interface {
	RefreshToken(ctx context.Context, req adminAuthModel.RefreshTokenRequest) (*adminAuthModel.RefreshTokenResponse, error)
}
