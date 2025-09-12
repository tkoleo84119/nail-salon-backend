package adminAuth

import (
	"context"
	"errors"
	"net/netip"
	"time"

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

func (s *RefreshToken) RefreshToken(ctx context.Context, req adminAuthModel.RefreshTokenRequest, refreshTokenCtx adminAuthModel.RefreshTokenContext) (*adminAuthModel.RefreshTokenResponse, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthStaffFailed)
		}
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

	// Rotate refresh token: revoke old, issue and store new
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// best-effort revoke old token; even if it fails due to race it won't block issuance
	_ = s.queries.RevokeStaffUserToken(ctx, req.RefreshToken)

	// store new refresh token with 7-day expiry (same policy as login)
	exp := time.Now().Add(7 * 24 * time.Hour)
	var ipAddr *netip.Addr
	if addr, err := netip.ParseAddr(refreshTokenCtx.IPAddress); err == nil {
		ipAddr = &addr
	}

	_, err = s.queries.CreateStaffUserToken(ctx, dbgen.CreateStaffUserTokenParams{
		ID:           utils.GenerateID(),
		StaffUserID:  staffUser.ID,
		RefreshToken: newRefreshToken,
		UserAgent:    utils.StringPtrToPgText(&refreshTokenCtx.UserAgent, true),
		IpAddress:    ipAddr,
		ExpiredAt:    utils.TimePtrToPgTimestamptz(&exp),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	// Build response with new refresh token (handler will set cookie, not JSON)
	return &adminAuthModel.RefreshTokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    s.jwtConfig.ExpiryHours * 3600,
		RefreshToken: newRefreshToken,
	}, nil
}
