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

// MockUpdateStaffService implements the UpdateStaffServiceInterface for testing
type MockUpdateStaffService struct {
	mock.Mock
}

// Ensure MockUpdateStaffService implements the interface
var _ staffService.UpdateStaffServiceInterface = (*MockUpdateStaffService)(nil)

func (m *MockUpdateStaffService) UpdateStaff(ctx context.Context, targetID string, req staff.UpdateStaffRequest, updaterID int64, updaterRole string) (*staff.UpdateStaffResponse, error) {
	args := m.Called(ctx, targetID, req, updaterID, updaterRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.UpdateStaffResponse), args.Error(1)
}

func setupTestGinForUpdate() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestUpdateStaffHandler_UpdateStaff_Success(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data
	role := staff.RoleManager
	isActive := false
	req := staff.UpdateStaffRequest{
		Role:     &role,
		IsActive: &isActive,
	}

	expectedResponse := &staff.UpdateStaffResponse{
		ID:       "123456789",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: false,
	}

	// Set up mock expectations
	mockService.On("UpdateStaff", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin).Return(expectedResponse, nil)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

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
	assert.Equal(t, expectedResponse.IsActive, responseData["isActive"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateStaffHandler_UpdateStaff_MissingID(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context but no ID param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	// No ID param set

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "id為必填項目", response.Errors["id"])
}

func TestUpdateStaffHandler_UpdateStaff_MissingStaffContext(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Call handler
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestUpdateStaffHandler_UpdateStaff_InvalidJSON(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestUpdateStaffHandler_UpdateStaff_EmptyRequest(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Empty request
	req := staff.UpdateStaffRequest{}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "至少需要提供一個欄位進行更新", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "至少需要提供一個欄位進行更新", response.Errors["request"])
}

func TestUpdateStaffHandler_UpdateStaff_InvalidRole(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data with invalid role
	invalidRole := "INVALID_ROLE"
	req := staff.UpdateStaffRequest{
		Role: &invalidRole,
	}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "role只可以傳入特定值", response.Errors["role"])
}

func TestUpdateStaffHandler_UpdateStaff_ServiceError(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Set up mock expectations - return service error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("UpdateStaff", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin).Return(nil, serviceError)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

	// Assert response
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "權限不足，無法執行此操作", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateStaffHandler_UpdateStaff_InternalError(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateStaffService)
	handler := NewUpdateStaffHandler(mockService)

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("UpdateStaff", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin).Return(nil, internalError)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/staff/123456789", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

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
	handler.UpdateStaff(c)

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
