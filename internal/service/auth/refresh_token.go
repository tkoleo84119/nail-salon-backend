package auth

import (
	"context"
	"errors"
	"net/netip"
	"time"

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

func (s *RefreshToken) RefreshToken(ctx context.Context, req auth.RefreshTokenRequest, refreshTokenCtx auth.RefreshTokenContext) (*auth.RefreshTokenResponse, error) {
	// Validate refresh token exists and is not revoked/expired
	tokenRecord, err := s.queries.GetValidCustomerToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthRefreshTokenInvalid)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to validate refresh token", err)
	}

	var ipAddr *netip.Addr
	if refreshTokenCtx.IPAddress != "" {
		if addr, err := netip.ParseAddr(refreshTokenCtx.IPAddress); err == nil {
			ipAddr = &addr
		}
	}

	// Generate new access token
	accessToken, err := utils.GenerateCustomerJWT(s.jwtConfig, tokenRecord.CustomerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to generate access token", err)
	}

	// Rotate refresh token
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// best-effort revoke old token; ignore error to not block issuance
	_ = s.queries.RevokeCustomerToken(ctx, req.RefreshToken)

	exp := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.queries.CreateCustomerToken(ctx, dbgen.CreateCustomerTokenParams{
		ID:           utils.GenerateID(),
		CustomerID:   tokenRecord.CustomerID,
		RefreshToken: newRefreshToken,
		UserAgent:    utils.StringPtrToPgText(&refreshTokenCtx.UserAgent, true),
		IpAddress:    ipAddr,
		ExpiredAt:    utils.TimePtrToPgTimestamptz(&exp),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	// Build response
	return &auth.RefreshTokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    int(s.jwtConfig.ExpiryHours * 3600),
		RefreshToken: newRefreshToken,
	}, nil
}
