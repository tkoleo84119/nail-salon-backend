package auth

import (
	"context"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CustomerLineLoginService struct {
	queries       dbgen.Querier
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
}

func NewCustomerLineLoginService(queries dbgen.Querier, lineConfig config.LineConfig, jwtConfig config.JWTConfig) *CustomerLineLoginService {
	lineValidator := utils.NewLineValidator(lineConfig.ChannelID)
	return &CustomerLineLoginService{
		queries:       queries,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
	}
}

func (s *CustomerLineLoginService) CustomerLineLogin(ctx context.Context, req auth.CustomerLineLoginRequest, loginCtx auth.CustomerLoginContext) (*auth.CustomerLineLoginResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	customerAuth, err := s.queries.GetCustomerAuthByProviderUid(ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    auth.ProviderLine,
		ProviderUid: profile.ProviderUid,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			// Customer not registered, return registration info
			response := &auth.CustomerLineLoginResponse{
				NeedRegister: true,
				LineProfile: &auth.CustomerLineProfile{
					ProviderUid: profile.ProviderUid,
					Name:        profile.Name,
					Email:       profile.Email,
				},
			}
			return response, nil
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer auth", err)
	}

	// Customer exists, generate tokens
	accessToken, err := s.generateAccessToken(customerAuth.CustomerID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, customerAuth.CustomerID, loginCtx)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &auth.CustomerLineLoginResponse{
		NeedRegister: false,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		Customer: &auth.Customer{
			ID:    utils.FormatID(customerAuth.CustomerID),
			Name:  customerAuth.CustomerName,
			Phone: customerAuth.CustomerPhone,
		},
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *CustomerLineLoginService) generateAccessToken(customerID int64) (string, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *CustomerLineLoginService) generateRefreshToken(ctx context.Context, customerID int64, loginCtx auth.CustomerLoginContext) (string, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Generate snowflake ID for token record
	tokenID := utils.GenerateID()

	// Store refresh token in database
	expiresAt := utils.TimeToPgTimez(time.Now().Add(7 * 24 * time.Hour)) // 7 days
	userAgent := utils.StringToText(&loginCtx.UserAgent)

	var ipAddress *netip.Addr
	if loginCtx.IPAddress != "" {
		if addr, err := netip.ParseAddr(loginCtx.IPAddress); err == nil {
			ipAddress = &addr
		}
	}

	_, err = s.queries.CreateCustomerToken(ctx, dbgen.CreateCustomerTokenParams{
		ID:           tokenID,
		CustomerID:   customerID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IpAddress:    ipAddress,
		ExpiredAt:    expiresAt,
	})

	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	return refreshToken, nil
}
