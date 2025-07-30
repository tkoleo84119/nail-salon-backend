package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CustomerLineRegisterService struct {
	queries       dbgen.Querier
	db            *pgxpool.Pool
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
}

func NewCustomerLineRegisterService(queries dbgen.Querier, db *pgxpool.Pool, lineConfig config.LineConfig, jwtConfig config.JWTConfig) *CustomerLineRegisterService {
	lineValidator := utils.NewLineValidator(lineConfig.ChannelID)
	return &CustomerLineRegisterService{
		queries:       queries,
		db:            db,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
	}
}

func (s *CustomerLineRegisterService) CustomerLineRegister(ctx context.Context, req auth.CustomerLineRegisterRequest, loginCtx auth.CustomerLoginContext) (*auth.CustomerLineRegisterResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	_, err = s.queries.GetCustomerAuthByProviderUid(ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    auth.ProviderLine,
		ProviderUid: profile.ProviderUid,
	})

	if err == nil {
		// Customer already exists
		return nil, errorCodes.NewServiceError(errorCodes.CustomerAlreadyExists, "this line account has been registered", nil)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// Database error
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer existence", err)
	}

	// Prepare customer record
	customerID := utils.GenerateID()

	birthday, err := utils.DateStringToPgDate(req.Birthday)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid birthday format", err)
	}
	city := utils.StringPtrToPgText(&req.City, false)
	referrer := utils.StringPtrToPgText(&req.Referrer, false)
	customerNote := utils.StringPtrToPgText(&req.CustomerNote, false)
	isIntrovert := utils.BoolPtrToPgBool(req.IsIntrovert)

	// Convert favorite arrays to PostgreSQL text arrays
	var favoriteShapes, favoriteColors, favoriteStyles, referralSource []string
	if req.FavoriteShapes != nil {
		favoriteShapes = req.FavoriteShapes
	}
	if req.FavoriteColors != nil {
		favoriteColors = req.FavoriteColors
	}
	if req.FavoriteStyles != nil {
		favoriteStyles = req.FavoriteStyles
	}
	if req.ReferralSource != nil {
		referralSource = req.ReferralSource
	}

	// Create customer auth record
	authID := utils.GenerateID()
	var otherInfo []byte
	if profile.Email != nil {
		otherInfoMap := map[string]interface{}{
			"name":  profile.Name,
			"email": *profile.Email,
		}
		otherInfo, _ = json.Marshal(otherInfoMap)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	createdCustomer, err := qtx.CreateCustomer(ctx, dbgen.CreateCustomerParams{
		ID:             customerID,
		Name:           req.Name,
		Phone:          req.Phone,
		Birthday:       birthday,
		City:           city,
		FavoriteShapes: favoriteShapes,
		FavoriteColors: favoriteColors,
		FavoriteStyles: favoriteStyles,
		IsIntrovert:    isIntrovert,
		ReferralSource: referralSource,
		Referrer:       referrer,
		CustomerNote:   customerNote,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer", err)
	}

	_, err = qtx.CreateCustomerAuth(ctx, dbgen.CreateCustomerAuthParams{
		ID:          authID,
		CustomerID:  customerID,
		Provider:    auth.ProviderLine,
		ProviderUid: profile.ProviderUid,
		OtherInfo:   otherInfo,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer auth", err)
	}

	// Generate access token
	accessToken, err := s.generateAccessToken(customerID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(ctx, qtx, customerID, loginCtx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Build response
	response := &auth.CustomerLineRegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Customer: &auth.CustomerRegisteredCustomer{
			ID:       utils.FormatID(createdCustomer.ID),
			Name:     createdCustomer.Name,
			Phone:    createdCustomer.Phone,
			Birthday: createdCustomer.Birthday.Time.Format("2006-01-02"),
		},
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *CustomerLineRegisterService) generateAccessToken(customerID int64) (string, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *CustomerLineRegisterService) generateRefreshToken(ctx context.Context, queries dbgen.Querier, customerID int64, loginCtx auth.CustomerLoginContext) (string, error) {
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

	_, err = queries.CreateCustomerToken(ctx, dbgen.CreateCustomerTokenParams{
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
