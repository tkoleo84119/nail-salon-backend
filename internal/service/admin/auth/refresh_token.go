package adminAuth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshToken struct {
	queries   *dbgen.Queries
	jwtConfig config.JWTConfig
}

func NewRefreshToken(queries *dbgen.Queries, jwtConfig config.JWTConfig) RefreshTokenInterface {
	return &RefreshToken{
		queries:   queries,
		jwtConfig: jwtConfig,
	}
}

func (s *RefreshToken) RefreshToken(ctx context.Context, req adminAuthModel.RefreshTokenRequest) (*adminAuthModel.RefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	refreshTokenInfo, err := s.queries.GetValidStaffUserToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to validate refresh token", err)
	}

	// Get staff user information to rebuild JWT claims
	staffUser, err := s.queries.GetStaffUserByID(ctx, refreshTokenInfo.StaffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get staff user", err)
	}
	if !staffUser.IsActive.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthStaffFailed)
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(s.jwtConfig, staffUser.ID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	// Build response
	return &adminAuthModel.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.jwtConfig.ExpiryHours * 3600,
	}, nil
}
