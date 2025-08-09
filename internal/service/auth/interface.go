package auth

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
)

type LineRegisterInterface interface {
	LineRegister(ctx context.Context, req auth.LineRegisterRequest, loginCtx auth.LoginContext) (*auth.LineRegisterResponse, error)
}

type LineLoginInterface interface {
	LineLogin(ctx context.Context, req auth.LineLoginRequest, loginCtx auth.LoginContext) (*auth.LineLoginResponse, error)
}

type RefreshTokenInterface interface {
	RefreshToken(ctx context.Context, req auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error)
}
