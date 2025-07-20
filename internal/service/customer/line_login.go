package customer

import (
	"context"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineLoginService struct {
	queries       dbgen.Querier
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
}

func NewLineLoginService(queries dbgen.Querier, lineConfig config.LineConfig, jwtConfig config.JWTConfig) *LineLoginService {
	lineValidator := utils.NewLineValidator(lineConfig.ChannelID)
	return &LineLoginService{
		queries:       queries,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
	}
}

func (s *LineLoginService) LineLogin(ctx context.Context, req customer.LineLoginRequest, loginCtx customer.LoginContext) (*customer.LineLoginResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	customerAuth, err := s.queries.GetCustomerAuthByProviderUid(ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    customer.ProviderLine,
		ProviderUid: profile.ProviderUid,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			// Customer not registered, return registration info
			response := &customer.LineLoginResponse{
				NeedRegister: true,
				LineProfile: &customer.LineProfile{
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
	response := &customer.LineLoginResponse{
		NeedRegister: false,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		Customer: &customer.Customer{
			ID:    utils.FormatID(customerAuth.CustomerID),
			Name:  customerAuth.CustomerName,
			Phone: customerAuth.CustomerPhone,
		},
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *LineLoginService) generateAccessToken(customerID int64) (string, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *LineLoginService) generateRefreshToken(ctx context.Context, customerID int64, loginCtx customer.LoginContext) (string, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Generate snowflake ID for token record
	tokenID := utils.GenerateID()

	// Store refresh token in database
	expiresAt := pgtype.Timestamptz{
		Time:  time.Now().Add(30 * 24 * time.Hour), // 7 days
		Valid: true,
	}

	var ipAddress *netip.Addr
	if loginCtx.IPAddress != "" {
		if addr, err := netip.ParseAddr(loginCtx.IPAddress); err == nil {
			ipAddress = &addr
		}
	}

	var userAgent pgtype.Text
	if loginCtx.UserAgent != "" {
		userAgent = pgtype.Text{String: loginCtx.UserAgent, Valid: true}
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