package customer

import (
	"context"
	"encoding/json"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineRegisterService struct {
	queries       dbgen.Querier
	db            *pgxpool.Pool
	lineValidator *utils.LineValidator
	jwtConfig     config.JWTConfig
}

func NewLineRegisterService(queries dbgen.Querier, db *pgxpool.Pool, lineConfig config.LineConfig, jwtConfig config.JWTConfig) *LineRegisterService {
	lineValidator := utils.NewLineValidator(lineConfig.ChannelID)
	return &LineRegisterService{
		queries:       queries,
		db:            db,
		lineValidator: lineValidator,
		jwtConfig:     jwtConfig,
	}
}

func (s *LineRegisterService) LineRegister(ctx context.Context, req customer.LineRegisterRequest, loginCtx customer.LoginContext) (*customer.LineRegisterResponse, error) {
	// Validate LINE ID token and get profile
	profile, err := s.lineValidator.ValidateIdToken(req.IdToken)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	_, err = s.queries.GetCustomerAuthByProviderUid(ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    customer.ProviderLine,
		ProviderUid: profile.ProviderUid,
	})

	if err == nil {
		// Customer already exists
		return nil, errorCodes.NewServiceError(errorCodes.CustomerAlreadyExists, "this line account has been registered", nil)
	} else if err != pgx.ErrNoRows {
		// Database error
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer existence", err)
	}

	// Create customer record
	customerID := utils.GenerateID()

	// Parse and convert birthday to pgtype.Date
	parsedBirthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid birthday format", err)
	}

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

	// Prepare optional fields
	var city pgtype.Text
	if req.City != "" {
		city = pgtype.Text{String: req.City, Valid: true}
	}

	var referrer pgtype.Text
	if req.Referrer != "" {
		referrer = pgtype.Text{String: req.Referrer, Valid: true}
	}

	var customerNote pgtype.Text
	if req.CustomerNote != "" {
		customerNote = pgtype.Text{String: req.CustomerNote, Valid: true}
	}

	var isIntrovert pgtype.Bool
	if req.IsIntrovert != nil {
		isIntrovert = pgtype.Bool{Bool: *req.IsIntrovert, Valid: true}
	}

	birthday := pgtype.Date{
		Time:  parsedBirthday,
		Valid: true,
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

	_, err = qtx.CreateCustomerAuth(ctx, dbgen.CreateCustomerAuthParams{
		ID:          authID,
		CustomerID:  customerID,
		Provider:    customer.ProviderLine,
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
	response := &customer.LineRegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Customer: &customer.RegisteredCustomer{
			ID:       utils.FormatID(createdCustomer.ID),
			Name:     createdCustomer.Name,
			Phone:    createdCustomer.Phone,
			Birthday: createdCustomer.Birthday.Time.Format("2006-01-02"),
		},
	}

	return response, nil
}

// generateAccessToken generates a JWT access token for the customer
func (s *LineRegisterService) generateAccessToken(customerID int64) (string, error) {
	token, err := utils.GenerateCustomerJWT(s.jwtConfig, customerID)
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate access token", err)
	}

	return token, nil
}

// generateRefreshToken generates and stores a refresh token for the customer
func (s *LineRegisterService) generateRefreshToken(ctx context.Context, queries dbgen.Querier, customerID int64, loginCtx customer.LoginContext) (string, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to generate refresh token", err)
	}

	// Generate snowflake ID for token record
	tokenID := utils.GenerateID()

	// Store refresh token in database
	expiresAt := pgtype.Timestamptz{
		Time:  time.Now().Add(30 * 24 * time.Hour), // 30 days
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
