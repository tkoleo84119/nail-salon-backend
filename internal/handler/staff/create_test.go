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

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
)

// MockCreateStaffService implements the CreateStaffServiceInterface for testing
type MockCreateStaffService struct {
	mock.Mock
}

// Ensure MockCreateStaffService implements the interface
var _ staffService.CreateStaffServiceInterface = (*MockCreateStaffService)(nil)

func (m *MockCreateStaffService) CreateStaff(ctx context.Context, req staff.CreateStaffRequest, creatorRole string, creatorStoreIDs []int64) (*staff.CreateStaffResponse, error) {
	args := m.Called(ctx, req, creatorRole, creatorStoreIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.CreateStaffResponse), args.Error(1)
}

func setupTestGinForCreate() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestCreateStaffHandler_CreateStaff_Success(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	// Set up mock expectations
	expectedResponse := &staff.CreateStaffResponse{
		ID:       "123456789",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}

	mockService.On("CreateStaff", mock.Anything, mock.AnythingOfType("staff.CreateStaffRequest"), "SUPER_ADMIN", []int64{1}).Return(expectedResponse, nil)

	// Create request
	createReq := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify the response data
	responseData := response.Data.(map[string]interface{})
	assert.Equal(t, expectedResponse.ID, responseData["id"])
	assert.Equal(t, expectedResponse.Username, responseData["username"])
	assert.Equal(t, expectedResponse.Email, responseData["email"])
	assert.Equal(t, expectedResponse.Role, responseData["role"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateStaffHandler_CreateStaff_ValidationError(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	tests := []struct {
		name     string
		request  staff.CreateStaffRequest
		expected string
	}{
		{
			name: "missing_username",
			request: staff.CreateStaffRequest{
				Email:    "test@example.com",
				Password: "testpassword",
				Role:     staff.RoleManager,
				StoreIDs: []string{"1"},
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "missing_email",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Password: "testpassword",
				Role:     staff.RoleManager,
				StoreIDs: []string{"1"},
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "missing_password",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Role:     staff.RoleManager,
				StoreIDs: []string{"1"},
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "missing_role",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
				StoreIDs: []string{"1"},
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "missing_store_ids",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
				Role:     staff.RoleManager,
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "invalid_email",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "testpassword",
				Role:     staff.RoleManager,
				StoreIDs: []string{"1"},
			},
			expected: "輸入驗證失敗",
		},
		{
			name: "empty_store_ids",
			request: staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
				Role:     staff.RoleManager,
				StoreIDs: []string{},
			},
			expected: "輸入驗證失敗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			reqBody, _ := json.Marshal(tt.request)

			// Create test context with staff context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set staff context
			staffContext := &common.StaffContext{
				UserID:   "1",
				Username: "admin",
				Role:     staff.RoleSuperAdmin,
				StoreList: []common.Store{
					{ID: "1", Name: "Test Store"},
				},
			}
			c.Set("user", staffContext)

			// Call handler
			handler.CreateStaff(c)

			// Assert response
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response common.ApiResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, response.Message)
			assert.Nil(t, response.Data)
			assert.NotNil(t, response.Errors)
		})
	}
}

func TestCreateStaffHandler_CreateStaff_InvalidJSON(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "JSON格式錯誤", response.Errors["request"])
}

func TestCreateStaffHandler_CreateStaff_MissingStaffContext(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	// Create request
	createReq := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	handler.CreateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestCreateStaffHandler_CreateStaff_ServiceError(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	// Set up mock expectations - return service error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.UserAlreadyExists)
	mockService.On("CreateStaff", mock.Anything, mock.AnythingOfType("staff.CreateStaffRequest"), "SUPER_ADMIN", []int64{1}).Return(nil, serviceError)

	// Create request
	createReq := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusConflict, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "帳號或Email已存在", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateStaffHandler_CreateStaff_InternalError(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateStaffService)
	handler := NewCreateStaffHandler(mockService)

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("CreateStaff", mock.Anything, mock.AnythingOfType("staff.CreateStaffRequest"), "SUPER_ADMIN", []int64{1}).Return(nil, internalError)

	// Create request
	createReq := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "系統發生錯誤，請稍後再試", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}