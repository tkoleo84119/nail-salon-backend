package customer

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

func init() {
	// Initialize snowflake for testing
	utils.InitSnowflake(1)
}

// getTestLineRegisterConfig returns test configs for registration
func getTestLineRegisterConfig() (config.LineConfig, config.JWTConfig) {
	lineConfig := config.LineConfig{
		ChannelID: "YOUR_LINE_CHANNEL_ID", // This triggers test mode
	}
	jwtConfig := config.JWTConfig{
		Secret:      "test-jwt-secret",
		ExpiryHours: 24,
	}
	return lineConfig, jwtConfig
}

func TestLineRegisterService_LineRegister_InvalidBirthday(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	lineConfig, jwtConfig := getTestLineRegisterConfig()
	service := NewLineRegisterService(mockQuerier, mockDB, lineConfig, jwtConfig)

	ctx := context.Background()
	req := customer.LineRegisterRequest{
		IdToken:  "valid-token-no-email.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "invalid-date", // Invalid birthday format
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock customer does not exist
	mockQuerier.On("GetCustomerAuthByProviderUid", mock.Anything, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    customer.ProviderLine,
		ProviderUid: "U12345678",
	}).Return(dbgen.GetCustomerAuthByProviderUidRow{}, pgx.ErrNoRows)

	response, err := service.LineRegister(ctx, req, loginCtx)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestLineRegisterService_LineRegister_CustomerAlreadyExists(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	lineConfig, jwtConfig := getTestLineRegisterConfig()
	service := NewLineRegisterService(mockQuerier, mockDB, lineConfig, jwtConfig)

	ctx := context.Background()
	req := customer.LineRegisterRequest{
		IdToken:  "valid-token-no-email.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "1990-01-01",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock customer already exists
	existingCustomer := dbgen.GetCustomerAuthByProviderUidRow{
		ID:            1,
		CustomerID:    1000000001,
		Provider:      customer.ProviderLine,
		ProviderUid:   "U12345678",
		CustomerName:  "小美",
		CustomerPhone: "09xxxxxxxx",
	}

	mockQuerier.On("GetCustomerAuthByProviderUid", ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    customer.ProviderLine,
		ProviderUid: "U12345678",
	}).Return(existingCustomer, nil)

	response, err := service.LineRegister(ctx, req, loginCtx)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.CustomerAlreadyExists, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestLineRegisterService_LineRegister_InvalidToken(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	lineConfig, jwtConfig := getTestLineRegisterConfig()
	service := NewLineRegisterService(mockQuerier, mockDB, lineConfig, jwtConfig)

	ctx := context.Background()
	req := customer.LineRegisterRequest{
		IdToken:  "invalid-token.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "1990-01-01",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	response, err := service.LineRegister(ctx, req, loginCtx)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthLineTokenInvalid, serviceErr.Code)

	// Verify no database calls were made
	mockQuerier.AssertNotCalled(t, "GetCustomerAuthByProviderUid")
}

func TestLineRegisterService_LineRegister_DatabaseError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	lineConfig, jwtConfig := getTestLineRegisterConfig()
	service := NewLineRegisterService(mockQuerier, mockDB, lineConfig, jwtConfig)

	ctx := context.Background()
	req := customer.LineRegisterRequest{
		IdToken:  "valid-token-no-email.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "1990-01-01",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock database error
	mockQuerier.On("GetCustomerAuthByProviderUid", ctx, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    customer.ProviderLine,
		ProviderUid: "U12345678",
	}).Return(dbgen.GetCustomerAuthByProviderUidRow{}, assert.AnError)

	response, err := service.LineRegister(ctx, req, loginCtx)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

// Test helper function for generating access tokens
func TestLineRegisterService_generateAccessToken_Success(t *testing.T) {
	_, jwtConfig := getTestLineRegisterConfig()
	service := &LineRegisterService{
		jwtConfig: jwtConfig,
	}

	customerID := int64(1000000001)
	token, err := service.generateAccessToken(customerID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLineRegisterService_generateAccessToken_InvalidCustomerID(t *testing.T) {
	_, jwtConfig := getTestLineRegisterConfig()
	service := &LineRegisterService{
		jwtConfig: jwtConfig,
	}

	// Test with zero customer ID which should still work
	customerID := int64(0)
	token, err := service.generateAccessToken(customerID)

	// JWT generation should still work with customer ID 0
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
