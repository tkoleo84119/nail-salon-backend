package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshTokenService struct {
	queries   *dbgen.Queries
	jwtConfig config.JWTConfig
}

func NewRefreshTokenService(queries *dbgen.Queries, jwtConfig config.JWTConfig) *RefreshTokenService {
	return &RefreshTokenService{
		queries:   queries,
		jwtConfig: jwtConfig,
	}
}

func (s *RefreshTokenService) RefreshToken(ctx context.Context, req auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	tokenRecord, err := s.queries.GetValidCustomerToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to validate refresh token", err)
	}

	// Generate new access token
	accessToken, err := utils.GenerateCustomerJWT(s.jwtConfig, tokenRecord.CustomerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to generate access token", err)
	}

	// Build response
	return &auth.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.jwtConfig.ExpiryHours * 3600,
	}, nil
}