package auth

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
)

type LoginServiceInterface interface {
	Login(ctx context.Context, req auth.LoginRequest, loginCtx auth.LoginContext) (*auth.LoginResponse, error)
}
