package auth

import (
	"context"
	"log"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineRegister struct {
	queries       *dbgen.Queries
	db            *pgxpool.Pool
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
	activityLog   cache.ActivityLogCacheInterface
}

func NewLineRegister(queries *dbgen.Queries, db *pgxpool.Pool, lineConfig config.LineConfig, jwtConfig config.JWTConfig, activityLog cache.ActivityLogCacheInterface) *LineRegister {
	lineValidator := utils.NewLineValidator(lineConfig.LiffChannelID)
	return &LineRegister{
		queries:       queries,
		db:            db,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
		activityLog:   activityLog,
	}
}

func (s *LineRegister) LineRegister(ctx context.Context, req auth.LineRegisterRequest, loginCtx auth.LoginContext) (*auth.LineRegisterResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	exist, err := s.queries.CheckCustomerExistsByLineUid(ctx, profile.ProviderUid)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer existence", err)
	}
	if exist {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerAlreadyExists)
	}

	// Prepare customer record
	customerID := utils.GenerateID()

	birthday, err := utils.DateStringToPgDate(req.Birthday)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValFieldDateFormat)
	}
	defaultLevel := "NORMAL"

	// Convert favorite arrays to PostgreSQL text arrays
	var favoriteShapes, favoriteColors, favoriteStyles, referralSource []string
	if req.FavoriteShapes != nil {
		favoriteShapes = *req.FavoriteShapes
	}
	if req.FavoriteColors != nil {
		favoriteColors = *req.FavoriteColors
	}
	if req.FavoriteStyles != nil {
		favoriteStyles = *req.FavoriteStyles
	}
	if req.ReferralSource != nil {
		referralSource = *req.ReferralSource
	}

	// start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	err = qtx.CreateCustomer(ctx, dbgen.CreateCustomerParams{
		ID:             customerID,
		LineUid:        profile.ProviderUid,
		LineName:       utils.StringPtrToPgText(&profile.Name, true),
		Email:          utils.StringPtrToPgText(req.Email, true),
		Name:           req.Name,
		Phone:          req.Phone,
		Birthday:       birthday,
		City:           utils.StringPtrToPgText(req.City, true),
		FavoriteShapes: favoriteShapes,
		FavoriteColors: favoriteColors,
		FavoriteStyles: favoriteStyles,
		IsIntrovert:    utils.BoolPtrToPgBool(req.IsIntrovert),
		ReferralSource: referralSource,
		Referrer:       utils.StringPtrToPgText(req.Referrer, true),
		CustomerNote:   utils.StringPtrToPgText(req.CustomerNote, true),
		Level:          utils.StringPtrToPgText(&defaultLevel, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer", err)
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

	// Generate customer terms acceptance
	customerTermsAcceptanceID := utils.GenerateID()
	acceptedAt := time.Now()
	acceptedAtPg := utils.TimePtrToPgTimestamptz(&acceptedAt)
	err = qtx.CreateCustomerTermsAcceptance(ctx, dbgen.CreateCustomerTermsAcceptanceParams{
		ID:           customerTermsAcceptanceID,
		CustomerID:   customerID,
		TermsVersion: common.CustomerTermsVersion1,
		AcceptedAt:   acceptedAtPg,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer terms acceptance", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to commit transaction", err)
	}

	// Log activity
	go func() {
		logCtx := context.Background()
		if err := s.activityLog.LogCustomerRegister(logCtx, req.Name); err != nil {
			log.Printf("failed to log customer register activity: %v", err)
		}
	}()

	// Build response
	response := &auth.LineRegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtConfig.ExpiryHours * 3600,
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *LineRegister) generateAccessToken(customerID int64) (string, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *LineRegister) generateRefreshToken(ctx context.Context, queries dbgen.Querier, customerID int64, loginCtx auth.LoginContext) (string, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Generate snowflake ID for token record
	tokenID := utils.GenerateID()

	// Store refresh token in database
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	expiresAtPg := utils.TimePtrToPgTimestamptz(&expiresAt)
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
		ExpiredAt:    expiresAtPg,
	})

	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to store refresh token", err)
	}

	return refreshToken, nil
}
