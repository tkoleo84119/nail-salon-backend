package customer

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

// getTestLineConfig returns a test LINE configuration
func getTestLineConfig() config.LineConfig {
	return config.LineConfig{
		ChannelID:        "YOUR_LINE_CHANNEL_ID", // This triggers test mode
	}
}

// getTestJWTConfig returns a test JWT configuration
func getTestJWTConfig() config.JWTConfig {
	return config.JWTConfig{
		Secret:      "test-jwt-secret",
		ExpiryHours: 24,
	}
}

func TestLineLoginService_LineLogin_CustomerNotRegistered(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "valid-token-no-email.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock customer auth not found
	mockQuerier.On("GetCustomerAuthByProviderUid", mock.Anything, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    "LINE",
		ProviderUid: "U12345678",
	}).Return(dbgen.GetCustomerAuthByProviderUidRow{}, pgx.ErrNoRows)

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.NeedRegister)
	assert.NotNil(t, response.LineProfile)
	assert.Equal(t, "U12345678", response.LineProfile.ProviderUid)
	assert.Equal(t, "Test User", response.LineProfile.Name)
	assert.Nil(t, response.LineProfile.Email)
	assert.Nil(t, response.AccessToken)
	assert.Nil(t, response.RefreshToken)
	assert.Nil(t, response.Customer)

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
}

func TestLineLoginService_LineLogin_CustomerNotRegistered_WithEmail(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request with email token
	req := customer.LineLoginRequest{
		IdToken: "valid-token-with-email.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock customer auth not found
	mockQuerier.On("GetCustomerAuthByProviderUid", mock.Anything, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    "LINE",
		ProviderUid: "U12345678",
	}).Return(dbgen.GetCustomerAuthByProviderUidRow{}, pgx.ErrNoRows)

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.NeedRegister)
	assert.NotNil(t, response.LineProfile)
	assert.Equal(t, "U12345678", response.LineProfile.ProviderUid)
	assert.Equal(t, "Test User", response.LineProfile.Name)
	assert.NotNil(t, response.LineProfile.Email)
	assert.Equal(t, "test@example.com", *response.LineProfile.Email)

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
}

func TestLineLoginService_LineLogin_CustomerAlreadyRegistered(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "valid-token-no-email.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock existing customer auth
	customerAuthRow := dbgen.GetCustomerAuthByProviderUidRow{
		ID:           1,
		CustomerID:   1000000001,
		Provider:     "LINE",
		ProviderUid:  "U12345678",
		CustomerName: "小美",
		CustomerPhone: "09xxxxxxxx",
	}

	mockQuerier.On("GetCustomerAuthByProviderUid", mock.Anything, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    "LINE",
		ProviderUid: "U12345678",
	}).Return(customerAuthRow, nil)

	// Mock token creation
	mockQuerier.On("CreateCustomerToken", mock.Anything, mock.MatchedBy(func(params dbgen.CreateCustomerTokenParams) bool {
		return params.CustomerID == 1000000001
	})).Return(dbgen.CustomerToken{
		ID:           2000000001,
		CustomerID:   1000000001,
		RefreshToken: "mock-refresh-token",
	}, nil)

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.NeedRegister)
	assert.NotNil(t, response.AccessToken)
	assert.NotNil(t, response.RefreshToken)
	assert.NotNil(t, response.Customer)
	assert.Equal(t, "1000000001", response.Customer.ID)
	assert.Equal(t, "小美", response.Customer.Name)
	assert.Equal(t, "09xxxxxxxx", response.Customer.Phone)
	assert.Nil(t, response.LineProfile)

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
}

func TestLineLoginService_LineLogin_InvalidToken(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request with invalid token
	req := customer.LineLoginRequest{
		IdToken: "invalid-token.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "AUTH_LINE_TOKEN_INVALID", serviceErr.Code)
	}

	// Verify no database calls were made
	mockQuerier.AssertNotCalled(t, "GetCustomerAuthByProviderUid")
}

func TestLineLoginService_LineLogin_ExpiredToken(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request with expired token
	req := customer.LineLoginRequest{
		IdToken: "expired-token.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "AUTH_LINE_TOKEN_EXPIRED", serviceErr.Code)
	}

	// Verify no database calls were made
	mockQuerier.AssertNotCalled(t, "GetCustomerAuthByProviderUid")
}

func TestLineLoginService_LineLogin_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "valid-token-no-email.payload.signature",
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Mock database error
	mockQuerier.On("GetCustomerAuthByProviderUid", mock.Anything, dbgen.GetCustomerAuthByProviderUidParams{
		Provider:    "LINE",
		ProviderUid: "U12345678",
	}).Return(dbgen.GetCustomerAuthByProviderUidRow{}, assert.AnError)

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "SYS_DATABASE_ERROR", serviceErr.Code)
	}

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
}

func TestLineLoginService_LineLogin_InvalidTokenStructure(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLineLoginService(mockQuerier, getTestLineConfig(), getTestJWTConfig())

	// Test request with invalid token structure (not 3 parts)
	req := customer.LineLoginRequest{
		IdToken: "invalid.token", // Only 2 parts instead of 3
	}

	loginCtx := customer.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call the service
	response, err := service.LineLogin(context.Background(), req, loginCtx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "AUTH_LINE_TOKEN_INVALID", serviceErr.Code)
		assert.Contains(t, serviceErr.Message, "invalid token structure")
	}

	// Verify no database calls were made
	mockQuerier.AssertNotCalled(t, "GetCustomerAuthByProviderUid")
}