package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
)

// MockLineLoginService implements the LineLoginServiceInterface for testing
type MockLineLoginService struct {
	mock.Mock
}

// Ensure MockLineLoginService implements the interface
var _ customerService.LineLoginServiceInterface = (*MockLineLoginService)(nil)

func (m *MockLineLoginService) LineLogin(ctx context.Context, req customer.LineLoginRequest, loginCtx customer.LoginContext) (*customer.LineLoginResponse, error) {
	args := m.Called(ctx, req, loginCtx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.LineLoginResponse), args.Error(1)
}

func setupTestGinForLineLogin() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestLineLoginHandler_LineLogin_Success_NeedRegister(t *testing.T) {
	setupTestGinForLineLogin()

	// Create mock service
	mockService := new(MockLineLoginService)
	handler := NewLineLoginHandler(mockService)

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "valid-token",
	}

	// Expected response from service
	expectedResponse := &customer.LineLoginResponse{
		NeedRegister: true,
		LineProfile: &customer.LineProfile{
			ProviderUid: "U12345678",
			Name:        "Mei",
			Email:       func() *string { s := "mei@example.com"; return &s }(),
		},
	}

	// Setup mock expectations - using mock.MatchedBy for LoginContext comparison
	mockService.On("LineLogin", mock.Anything, req, mock.MatchedBy(func(ctx customer.LoginContext) bool {
		return ctx.UserAgent == "test-agent" && ctx.IPAddress == "127.0.0.1"
	})).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/auth/customer/line/login", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "test-agent")
	httpReq.RemoteAddr = "127.0.0.1:12345"

	// Create gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call the handler
	handler.LineLogin(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Verify the response data
	responseData, _ := json.Marshal(response.Data)
	var actualResponse customer.LineLoginResponse
	json.Unmarshal(responseData, &actualResponse)

	assert.True(t, actualResponse.NeedRegister)
	assert.NotNil(t, actualResponse.LineProfile)
	assert.Equal(t, "U12345678", actualResponse.LineProfile.ProviderUid)
	assert.Equal(t, "Mei", actualResponse.LineProfile.Name)
	assert.Equal(t, "mei@example.com", *actualResponse.LineProfile.Email)
	assert.Nil(t, actualResponse.AccessToken)
	assert.Nil(t, actualResponse.RefreshToken)
	assert.Nil(t, actualResponse.Customer)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestLineLoginHandler_LineLogin_Success_AlreadyRegistered(t *testing.T) {
	setupTestGinForLineLogin()

	// Create mock service
	mockService := new(MockLineLoginService)
	handler := NewLineLoginHandler(mockService)

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "valid-token",
	}

	// Expected response from service
	accessToken := "access-token-123"
	refreshToken := "refresh-token-456"
	expectedResponse := &customer.LineLoginResponse{
		NeedRegister: false,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		Customer: &customer.Customer{
			ID:    "1000000001",
			Name:  "小美",
			Phone: "09xxxxxxxx",
		},
	}

	// Setup mock expectations
	mockService.On("LineLogin", mock.Anything, req, mock.MatchedBy(func(ctx customer.LoginContext) bool {
		return ctx.UserAgent == "test-agent" && ctx.IPAddress == "127.0.0.1"
	})).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/auth/customer/line/login", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "test-agent")
	httpReq.RemoteAddr = "127.0.0.1:12345"

	// Create gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call the handler
	handler.LineLogin(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Verify the response data
	responseData, _ := json.Marshal(response.Data)
	var actualResponse customer.LineLoginResponse
	json.Unmarshal(responseData, &actualResponse)

	assert.False(t, actualResponse.NeedRegister)
	assert.Equal(t, "access-token-123", *actualResponse.AccessToken)
	assert.Equal(t, "refresh-token-456", *actualResponse.RefreshToken)
	assert.NotNil(t, actualResponse.Customer)
	assert.Equal(t, "1000000001", actualResponse.Customer.ID)
	assert.Equal(t, "小美", actualResponse.Customer.Name)
	assert.Equal(t, "09xxxxxxxx", actualResponse.Customer.Phone)
	assert.Nil(t, actualResponse.LineProfile)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestLineLoginHandler_LineLogin_ValidationErrors(t *testing.T) {
	setupTestGinForLineLogin()

	// Create mock service
	mockService := new(MockLineLoginService)
	handler := NewLineLoginHandler(mockService)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
	}{
		{
			name: "Missing idToken",
			request: map[string]interface{}{
				"notIdToken": "some-value",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty idToken",
			request: customer.LineLoginRequest{
				IdToken: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "IdToken too long",
			request: customer.LineLoginRequest{
				IdToken: string(make([]byte, 501)), // Exceeds max length of 500
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			jsonData, _ := json.Marshal(tt.request)
			httpReq, _ := http.NewRequest("POST", "/api/auth/customer/line/login", bytes.NewBuffer(jsonData))
			httpReq.Header.Set("Content-Type", "application/json")

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httpReq

			// Call the handler
			handler.LineLogin(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}

	// Verify mock was not called
	mockService.AssertNotCalled(t, "LineLogin")
}

func TestLineLoginHandler_LineLogin_ServiceError(t *testing.T) {
	setupTestGinForLineLogin()

	// Create mock service
	mockService := new(MockLineLoginService)
	handler := NewLineLoginHandler(mockService)

	// Test request
	req := customer.LineLoginRequest{
		IdToken: "invalid-token",
	}

	// Setup mock to return error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthLineTokenInvalid)
	mockService.On("LineLogin", mock.Anything, req, mock.MatchedBy(func(ctx customer.LoginContext) bool {
		return true
	})).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/auth/customer/line/login", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "test-agent")

	// Create gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call the handler
	handler.LineLogin(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code) // AUTH_LINE_TOKEN_INVALID maps to 400

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestLineLoginHandler_LineLogin_InvalidJSON(t *testing.T) {
	setupTestGinForLineLogin()

	// Create mock service
	mockService := new(MockLineLoginService)
	handler := NewLineLoginHandler(mockService)

	// Create HTTP request with invalid JSON
	httpReq, _ := http.NewRequest("POST", "/api/auth/customer/line/login", bytes.NewBuffer([]byte("{invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call the handler
	handler.LineLogin(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify mock was not called
	mockService.AssertNotCalled(t, "LineLogin")
}