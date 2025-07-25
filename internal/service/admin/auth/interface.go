package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
)

type StaffLoginServiceInterface interface {
	StaffLogin(ctx context.Context, req adminAuthModel.StaffLoginRequest, loginCtx adminAuthModel.StaffLoginContext) (*adminAuthModel.StaffLoginResponse, error)
}
