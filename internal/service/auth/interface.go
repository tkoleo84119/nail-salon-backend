package auth

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
)

type CustomerLineLoginServiceInterface interface {
	CustomerLineLogin(ctx context.Context, req auth.CustomerLineLoginRequest, loginCtx auth.CustomerLoginContext) (*auth.CustomerLineLoginResponse, error)
}

type CustomerLineRegisterServiceInterface interface {
	CustomerLineRegister(ctx context.Context, req auth.CustomerLineRegisterRequest, loginCtx auth.CustomerLoginContext) (*auth.CustomerLineRegisterResponse, error)
}
