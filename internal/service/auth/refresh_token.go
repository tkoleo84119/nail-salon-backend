package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
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

func (s *RefreshToken) RefreshToken(ctx context.Context, req auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	tokenRecord, err := s.queries.GetValidCustomerToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
