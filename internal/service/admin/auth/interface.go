package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
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

type PermissionInterface interface {
	Permission(ctx context.Context, staffContext *common.StaffContext) (*adminAuthModel.PermissionResponse, error)
}

type UpdatePasswordInterface interface {
	UpdatePassword(ctx context.Context, req adminAuthModel.UpdatePasswordParsedRequest, staffContext *common.StaffContext) (*adminAuthModel.UpdatePasswordResponse, error)
}
