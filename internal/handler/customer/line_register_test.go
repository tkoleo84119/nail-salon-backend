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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

// MockLineRegisterService implements the LineRegisterServiceInterface for testing
type MockLineRegisterService struct {
	mock.Mock
}

// Ensure MockLineRegisterService implements the interface
var _ customerService.LineRegisterServiceInterface = (*MockLineRegisterService)(nil)

func (m *MockLineRegisterService) LineRegister(ctx context.Context, req customer.LineRegisterRequest, loginCtx customer.LoginContext) (*customer.LineRegisterResponse, error) {
	args := m.Called(ctx, req, loginCtx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.LineRegisterResponse), args.Error(1)
}

func TestLineRegisterHandler_LineRegister_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	// Setup mock service
	mockService := new(MockLineRegisterService)
	handler := NewLineRegisterHandler(mockService)

	// Test request
	isIntrovert := true
	reqBody := customer.LineRegisterRequest{
		IdToken:        "valid-token-no-email.payload.signature",
		Name:           "小美",
		Phone:          "0912345678",
		Birthday:       "1990-01-01",
		City:           "台北市",
		FavoriteShapes: []string{"圓形", "方形"},
		FavoriteColors: []string{"黑色", "白色"},
		FavoriteStyles: []string{"自然", "韓式"},
		IsIntrovert:    &isIntrovert,
		ReferralSource: []string{"朋友介紹", "網路廣告"},
		Referrer:       "1000000001",
		CustomerNote:   "這是客戶的備註",
	}

	expectedResponse := &customer.LineRegisterResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Customer: &customer.RegisteredCustomer{
			ID:       "2000000001",
			Name:     "小美",
			Phone:    "0912345678",
			Birthday: "1990-01-01",
		},
	}

	// Mock service response
	mockService.On("LineRegister", mock.Anything, mock.MatchedBy(func(req customer.LineRegisterRequest) bool {
		return req.Name == "小美" && req.Phone == "0912345678"
	}), mock.Anything).Return(expectedResponse, nil)

	// Create request
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/auth/customer/line/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.LineRegister(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "access-token", data["accessToken"])
	assert.Equal(t, "refresh-token", data["refreshToken"])

	customerData := data["customer"].(map[string]interface{})
	assert.Equal(t, "2000000001", customerData["id"])
	assert.Equal(t, "小美", customerData["name"])
	assert.Equal(t, "0912345678", customerData["phone"])
	assert.Equal(t, "1990-01-01", customerData["birthday"])

	mockService.AssertExpectations(t)
}

func TestLineRegisterHandler_LineRegister_ValidationError(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	// Setup mock service
	mockService := new(MockLineRegisterService)
	handler := NewLineRegisterHandler(mockService)

	tests := []struct {
		name        string
		requestBody map[string]interface{}
		expectedMsg string
	}{
		{
			name: "Missing required field - name",
			requestBody: map[string]interface{}{
				"idToken":  "valid-token.payload.signature",
				"phone":    "0912345678",
				"birthday": "1990-01-01",
			},
			expectedMsg: "name為必填項目",
		},
		{
			name: "Invalid phone format",
			requestBody: map[string]interface{}{
				"idToken":  "valid-token.payload.signature",
				"name":     "小美",
				"phone":    "0812345678", // Invalid - should start with 09
				"birthday": "1990-01-01",
			},
			expectedMsg: "phone必須為有效的台灣手機號碼格式 (例: 09xxxxxxxx)",
		},
		{
			name: "Missing required field - phone",
			requestBody: map[string]interface{}{
				"idToken":  "valid-token.payload.signature",
				"name":     "小美",
				"birthday": "1990-01-01",
			},
			expectedMsg: "phone為必填項目",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/auth/customer/line/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.LineRegister(c)

			// Assertions
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if errors, ok := response["errors"].(map[string]interface{}); ok {
				found := false
				for _, msg := range errors {
					if msg == tt.expectedMsg {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error message '%s' not found in errors: %v", tt.expectedMsg, errors)
			} else {
				t.Errorf("Expected errors field in response, got: %v", response)
			}
		})
	}

	mockService.AssertExpectations(t)
}

func TestLineRegisterHandler_LineRegister_CustomerAlreadyExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	// Setup mock service
	mockService := new(MockLineRegisterService)
	handler := NewLineRegisterHandler(mockService)

	// Test request
	reqBody := customer.LineRegisterRequest{
		IdToken:  "valid-token-no-email.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "1990-01-01",
	}

	// Mock service error - customer already exists
	mockService.On("LineRegister", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errorCodes.NewServiceError(errorCodes.CustomerAlreadyExists, "this line account has been registered", nil))

	// Create request
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/auth/customer/line/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.LineRegister(c)

	// Assertions
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["message"], "客戶已存在")

	mockService.AssertExpectations(t)
}

func TestLineRegisterHandler_LineRegister_InvalidToken(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	// Setup mock service
	mockService := new(MockLineRegisterService)
	handler := NewLineRegisterHandler(mockService)

	// Test request
	reqBody := customer.LineRegisterRequest{
		IdToken:  "invalid-token.payload.signature",
		Name:     "小美",
		Phone:    "0912345678",
		Birthday: "1990-01-01",
	}

	// Mock service error - invalid token
	mockService.On("LineRegister", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errorCodes.NewServiceError(errorCodes.AuthLineTokenInvalid, "token validation failed", nil))

	// Create request
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/auth/customer/line/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.LineRegister(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

func TestLineRegisterHandler_LineRegister_InvalidJSON(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	// Setup mock service
	mockService := new(MockLineRegisterService)
	handler := NewLineRegisterHandler(mockService)

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/auth/customer/line/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.LineRegister(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// The JSON parsing error returns a simple message format
	assert.Contains(t, response["message"], "輸入驗證失敗")

	mockService.AssertNotCalled(t, "LineRegister")
}
