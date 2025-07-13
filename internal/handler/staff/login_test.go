package staff

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// MockLoginService implements the LoginServiceInterface for testing
type MockLoginService struct {
	mock.Mock
}

// Ensure MockLoginService implements the interface
var _ staffService.LoginServiceInterface = (*MockLoginService)(nil)

func (m *MockLoginService) Login(ctx context.Context, req staff.LoginRequest, loginCtx staff.LoginContext) (*staff.LoginResponse, error) {
	args := m.Called(ctx, req, loginCtx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.LoginResponse), args.Error(1)
}

func setupTestGin() {
	gin.SetMode(gin.TestMode)
}

func TestLoginHandler_Login_Success(t *testing.T) {
	setupTestGin()

	// Create mock service
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	// Set up mock expectations
	expectedResponse := &staff.LoginResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    3600,
		User: staff.User{
			ID:       "123",
			Username: "testuser",
			Role:     "ADMIN",
			StoreList: []utils.Store{
				{ID: 1, Name: "Store 1"},
				{ID: 2, Name: "Store 2"},
			},
		},
	}

	mockService.On("Login", mock.Anything, mock.AnythingOfType("staff.LoginRequest"), mock.AnythingOfType("staff.LoginContext")).Return(expectedResponse, nil)

	// Create request
	loginReq := staff.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	reqBody, _ := json.Marshal(loginReq)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "test-agent")

	// Call handler
	handler.Login(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check that data is present and message is empty
	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Parse the data field
	dataBytes, _ := json.Marshal(response.Data)
	var loginResponse staff.LoginResponse
	err = json.Unmarshal(dataBytes, &loginResponse)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse.AccessToken, loginResponse.AccessToken)
	assert.Equal(t, expectedResponse.RefreshToken, loginResponse.RefreshToken)
	assert.Equal(t, expectedResponse.ExpiresIn, loginResponse.ExpiresIn)
	assert.Equal(t, expectedResponse.User, loginResponse.User)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestLoginHandler_Login_InvalidCredentials(t *testing.T) {
	setupTestGin()

	// Create mock service
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	// Set up mock expectations - invalid credentials
	mockService.On("Login", mock.Anything, mock.AnythingOfType("staff.LoginRequest"), mock.AnythingOfType("staff.LoginContext")).Return(nil, errors.New("invalid credentials"))

	// Create request
	loginReq := staff.LoginRequest{
		Username: "wronguser",
		Password: "wrongpassword",
	}
	reqBody, _ := json.Marshal(loginReq)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.Login(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "認證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "帳號或密碼錯誤", response.Errors["credentials"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestLoginHandler_Login_InternalError(t *testing.T) {
	setupTestGin()

	// Create mock service
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	// Set up mock expectations - internal error
	mockService.On("Login", mock.Anything, mock.AnythingOfType("staff.LoginRequest"), mock.AnythingOfType("staff.LoginContext")).Return(nil, errors.New("database connection failed"))

	// Create request
	loginReq := staff.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	reqBody, _ := json.Marshal(loginReq)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.Login(c)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "系統錯誤", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "伺服器內部錯誤", response.Errors["server"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestLoginHandler_Login_InvalidRequest(t *testing.T) {
	setupTestGin()

	// Create mock service (won't be called)
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	tests := []struct {
		name        string
		requestBody string
		expectCode  int
	}{
		{
			name:        "missing username",
			requestBody: `{"password":"test"}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "missing password",
			requestBody: `{"username":"test"}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "invalid JSON",
			requestBody: `{invalid json}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "empty request",
			requestBody: `{}`,
			expectCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/api/staff/login", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			assert.Equal(t, tt.expectCode, w.Code)

			var response common.ApiResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			// Should have message and errors, no data
			assert.NotEmpty(t, response.Message)
			assert.Nil(t, response.Data)
			assert.NotNil(t, response.Errors)
		})
	}
}
