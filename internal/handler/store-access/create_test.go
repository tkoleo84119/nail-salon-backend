package storeAccess

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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store-access"
	storeAccessService "github.com/tkoleo84119/nail-salon-backend/internal/service/store-access"
)

// MockCreateStoreAccessService implements the CreateStoreAccessServiceInterface for testing
type MockCreateStoreAccessService struct {
	mock.Mock
}

// Ensure MockCreateStoreAccessService implements the interface
var _ storeAccessService.CreateStoreAccessServiceInterface = (*MockCreateStoreAccessService)(nil)

func (m *MockCreateStoreAccessService) CreateStoreAccess(ctx context.Context, targetID string, req storeAccess.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*storeAccess.CreateStoreAccessResponse, bool, error) {
	args := m.Called(ctx, targetID, req, creatorID, creatorRole, creatorStoreIDs)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).(*storeAccess.CreateStoreAccessResponse), args.Bool(1), args.Error(2)
}

func setupTestGinForStoreAccess() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestCreateStoreAccessHandler_CreateStoreAccess_Success_Created(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	expectedResponse := &storeAccess.CreateStoreAccessResponse{
		StaffUserID: "123456789",
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}

	// Set up mock expectations
	mockService.On("CreateStoreAccess", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2}).Return(expectedResponse, true, nil)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

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
	assert.Equal(t, expectedResponse.StaffUserID, responseData["staffUserId"])
	assert.NotNil(t, responseData["storeList"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateStoreAccessHandler_CreateStoreAccess_Success_AlreadyExists(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	expectedResponse := &storeAccess.CreateStoreAccessResponse{
		StaffUserID: "123456789",
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}

	// Set up mock expectations - already exists
	mockService.On("CreateStoreAccess", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2}).Return(expectedResponse, false, nil)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateStoreAccessHandler_CreateStoreAccess_MissingID(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context but no ID param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	// No ID param set

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "員工ID為必填項目", response.Errors["id"])
}

func TestCreateStoreAccessHandler_CreateStoreAccess_MissingStaffContext(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Call handler
	handler.CreateStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestCreateStoreAccessHandler_CreateStoreAccess_InvalidJSON(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestCreateStoreAccessHandler_CreateStoreAccess_MissingStoreID(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data with missing store_id
	reqBody := []byte(`{}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "門市ID為必填項目", response.Errors["storeId"])
}

func TestCreateStoreAccessHandler_CreateStoreAccess_ServiceError(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Set up mock expectations - return service error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("CreateStoreAccess", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2}).Return(nil, false, serviceError)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

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

func TestCreateStoreAccessHandler_CreateStoreAccess_InternalError(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockCreateStoreAccessService)
	handler := NewCreateStoreAccessHandler(mockService)

	// Test data
	req := storeAccess.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("CreateStoreAccess", mock.Anything, "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2}).Return(nil, false, internalError)

	// Create request
	reqBody, _ := json.Marshal(req)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/staff/123456789/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123456789"}}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleSuperAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.CreateStoreAccess(c)

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
