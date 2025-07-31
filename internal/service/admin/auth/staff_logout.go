package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

// StaffLogoutServiceInterface defines the interface for staff logout service
type StaffLogoutServiceInterface interface {
	StaffLogout(ctx context.Context, req adminAuthModel.StaffLogoutRequest) (*adminAuthModel.StaffLogoutResponse, error)
}

type StaffLogoutService struct {
	repo *sqlxRepo.Repositories
}

func NewStaffLogoutService(repo *sqlxRepo.Repositories) StaffLogoutServiceInterface {
	return &StaffLogoutService{
		repo: repo,
	}
}

// StaffLogout revokes the refresh token and always returns success
func (s *StaffLogoutService) StaffLogout(ctx context.Context, req adminAuthModel.StaffLogoutRequest) (*adminAuthModel.StaffLogoutResponse, error) {
	_ = s.repo.StaffUserTokens.Revoke(ctx, req.RefreshToken)

	return &adminAuthModel.StaffLogoutResponse{
		Success: true,
	}, nil
}
