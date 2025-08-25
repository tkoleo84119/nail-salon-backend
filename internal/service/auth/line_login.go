package auth

import (
	"context"
	"errors"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineLogin struct {
	queries       dbgen.Querier
	db            *pgxpool.Pool
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
}

func NewLineLogin(queries dbgen.Querier, db *pgxpool.Pool, lineConfig config.LineConfig, jwtConfig config.JWTConfig) *LineLogin {
	lineValidator := utils.NewLineValidator(lineConfig.ChannelID)
	return &LineLogin{
		queries:       queries,
		db:            db,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
	}
}

func (s *LineLogin) LineLogin(ctx context.Context, req auth.LineLoginRequest, loginCtx auth.LoginContext) (*auth.LineLoginResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	customer, err := s.queries.GetCustomerByLineUid(ctx, profile.ProviderUid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response := &auth.LineLoginResponse{
				NeedRegister: true,
				LineProfile: &common.LineProfile{
					ProviderUid: profile.ProviderUid,
					Name:        profile.Name,
					Email:       profile.Email,
				},
			}
			return response, nil
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer exists", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Customer exists, generate tokens
	accessToken, expiresIn, err := s.generateAccessToken(customer.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, qtx, customer.ID, loginCtx)
	if err != nil {
		return nil, err
	}

	// check if customer line_name is different from profile.Name
	if profile.Name != "" && utils.PgTextToString(customer.LineName) != profile.Name {
		err = qtx.UpdateCustomerLineName(ctx, dbgen.UpdateCustomerLineNameParams{
			ID:       customer.ID,
			LineName: utils.StringPtrToPgText(&profile.Name, true),
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer line name", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Build response
	response := &auth.LineLoginResponse{
		NeedRegister: false,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		ExpiresIn:    &expiresIn,
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *LineLogin) generateAccessToken(customerID int64) (string, int, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	expiresIn := s.jwtConfig.ExpiryHours * 3600

	if err != nil {
		return "", 0, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, expiresIn, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *LineLogin) generateRefreshToken(ctx context.Context, qtx dbgen.Querier, customerID int64, loginCtx auth.LoginContext) (string, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Generate snowflake ID for token record
	tokenID := utils.GenerateID()

	// Store refresh token in database
	expiresAt := utils.TimeToPgTimestamptz(time.Now().Add(7 * 24 * time.Hour)) // 7 days
	userAgent := utils.StringPtrToPgText(&loginCtx.UserAgent, false)

	var ipAddress *netip.Addr
	if loginCtx.IPAddress != "" {
		if addr, err := netip.ParseAddr(loginCtx.IPAddress); err == nil {
			ipAddress = &addr
		}
	}

	_, err = qtx.CreateCustomerToken(ctx, dbgen.CreateCustomerTokenParams{
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
