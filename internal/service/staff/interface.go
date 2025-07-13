package staff

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
)

type LoginServiceInterface interface {
	Login(ctx context.Context, req staff.LoginRequest, loginCtx staff.LoginContext) (*staff.LoginResponse, error)
}