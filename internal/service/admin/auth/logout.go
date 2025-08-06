package adminAuth

import (
	"context"

	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

type Logout struct {
	queries *dbgen.Queries
}

func NewLogout(queries *dbgen.Queries) LogoutInterface {
	return &Logout{
		queries: queries,
	}
}

// StaffLogout revokes the refresh token and always returns success
func (s *Logout) Logout(ctx context.Context, req adminAuthModel.LogoutRequest) (*adminAuthModel.LogoutResponse, error) {
	_ = s.queries.RevokeStaffUserToken(ctx, req.RefreshToken)

	return &adminAuthModel.LogoutResponse{
		Success: true,
	}, nil
}
