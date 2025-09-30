package adminAuth

import (
	"context"
	"net/netip"
	"time"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Login struct {
	queries      *dbgen.Queries
	jwtConfig    config.JWTConfig
	cookieConfig config.CookieConfig
}

func NewLogin(
	queries *dbgen.Queries,
	jwtConfig config.JWTConfig,
	cookieConfig config.CookieConfig,
) LoginInterface {
	return &Login{
		queries:      queries,
		jwtConfig:    jwtConfig,
		cookieConfig: cookieConfig,
	}
}

func (s *Login) Login(ctx context.Context, req adminAuthModel.LoginRequest, loginCtx adminAuthModel.LoginContext) (*adminAuthModel.LoginResponse, error) {
	staffUser, err := s.queries.GetActiveStaffUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}
	if !utils.CheckPassword(req.Password, staffUser.PasswordHash) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(s.jwtConfig, staffUser.ID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Store refresh token
	tokenInfo := adminAuthModel.LoginTokenInfo{
		StaffUserID:  staffUser.ID,
		RefreshToken: refreshToken,
		Context:      loginCtx,
		ExpiresAt:    time.Now().Add(time.Duration(s.cookieConfig.AdminRefreshMaxAgeDays) * 24 * time.Hour),
	}

	if err := s.storeRefreshToken(ctx, tokenInfo); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	response := &adminAuthModel.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwtConfig.ExpiryHours * 3600),
	}

	return response, nil
}

// storeRefreshToken stores the refresh token in database
func (s *Login) storeRefreshToken(ctx context.Context, tokenInfo adminAuthModel.LoginTokenInfo) error {
	var ipAddr *netip.Addr
	if addr, err := netip.ParseAddr(tokenInfo.Context.IPAddress); err == nil {
		ipAddr = &addr
	}

	_, err := s.queries.CreateStaffUserToken(ctx, dbgen.CreateStaffUserTokenParams{
		ID:           utils.GenerateID(),
		StaffUserID:  tokenInfo.StaffUserID,
		RefreshToken: tokenInfo.RefreshToken,
		UserAgent:    utils.StringPtrToPgText(&tokenInfo.Context.UserAgent, true),
		IpAddress:    ipAddr,
		ExpiredAt:    utils.TimePtrToPgTimestamptz(&tokenInfo.ExpiresAt),
	})

	return err
}
