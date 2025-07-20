package customer

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
)

// LineLoginServiceInterface defines the interface for customer LINE login service
type LineLoginServiceInterface interface {
	LineLogin(ctx context.Context, req customer.LineLoginRequest, loginCtx customer.LoginContext) (*customer.LineLoginResponse, error)
}