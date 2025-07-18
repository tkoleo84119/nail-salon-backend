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

// MockDeleteStoreAccessService implements the DeleteStoreAccessServiceInterface for testing
type MockDeleteStoreAccessService struct {
	mock.Mock
}

// Ensure MockDeleteStoreAccessService implements the interface
var _ staffService.DeleteStoreAccessServiceInterface = (*MockDeleteStoreAccessService)(nil)

func (m *MockDeleteStoreAccessService) DeleteStoreAccess(ctx context.Context, targetID string, req staff.DeleteStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*staff.DeleteStoreAccessResponse, error) {
	args := m.Called(ctx, targetID, req, creatorID, creatorRole, creatorStoreIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.DeleteStoreAccessResponse), args.Error(1)
}

func setupTestContextWithStaffForDelete(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Mock staff context
	staffContext := common.StaffContext{
		UserID:   "123456789",
		Username: "testadmin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
			{ID: "2", Name: "Store 2"},
		},
	}
	c.Set("user", staffContext)

	// Set path parameter
	c.Params = []gin.Param{{Key: "id", Value: "987654321"}}

	return c, w
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_Success(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2"},
	}
	reqBody, _ := json.Marshal(req)

	expectedResponse := &staff.DeleteStoreAccessResponse{
		StaffUserID: "987654321",
		StoreList: []common.Store{
			{ID: "3", Name: "Store 3"},
		},
	}

	// Set up mock expectations
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(expectedResponse, nil)

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

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
	var deleteResponse staff.DeleteStoreAccessResponse
	err = json.Unmarshal(dataBytes, &deleteResponse)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse.StaffUserID, deleteResponse.StaffUserID)
	assert.Equal(t, expectedResponse.StoreList, deleteResponse.StoreList)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_MissingStaffContext(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Create test context without staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/api/staff/123/store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	// Call handler
	handler.DeleteStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify service was not called
	mockService.AssertNotCalled(t, "DeleteStoreAccess")
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_MissingStaffID(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Create test context with empty staff ID
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/api/staff//store-access", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Mock staff context
	staffContext := common.StaffContext{
		UserID:   "123456789",
		Username: "testadmin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Store 1"},
		},
	}
	c.Set("user", staffContext)

	// Set empty path parameter
	c.Params = []gin.Param{{Key: "id", Value: ""}}

	// Call handler
	handler.DeleteStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "id為必填項目", response.Errors["id"])

	// Verify service was not called
	mockService.AssertNotCalled(t, "DeleteStoreAccess")
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_InvalidRequestBody(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	tests := []struct {
		name        string
		requestBody string
		expectCode  int
	}{
		{
			name:        "invalid JSON",
			requestBody: `{invalid json}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "missing storeIds",
			requestBody: `{}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "empty storeIds array",
			requestBody: `{"storeIds": []}`,
			expectCode:  http.StatusBadRequest,
		},
		{
			name:        "null storeIds",
			requestBody: `{"storeIds": null}`,
			expectCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/123/store-access", []byte(tt.requestBody))

			handler.DeleteStoreAccess(c)

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

	// Verify service was not called for any test
	mockService.AssertNotCalled(t, "DeleteStoreAccess")
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_ServiceError_PermissionDenied(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Set up mock expectations - permission denied
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(nil, serviceError)

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

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

func TestDeleteStoreAccessHandler_DeleteStoreAccess_ServiceError_StaffNotFound(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Set up mock expectations - staff not found
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.UserStaffNotFound)
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(nil, serviceError)

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "指定的員工不存在", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_ServiceError_CannotModifySelf(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Set up mock expectations - cannot modify self
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.UserNotUpdateSelf)
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(nil, serviceError)

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "不可更新自己的帳號", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_InternalError(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1"},
	}
	reqBody, _ := json.Marshal(req)

	// Set up mock expectations - internal error
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(nil, errors.New("database connection failed"))

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

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

func TestDeleteStoreAccessHandler_DeleteStoreAccess_InvalidStoreIDsInRequest(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	tests := []struct {
		name              string
		requestBody       string
		expectCode        int
		shouldCallService bool
	}{
		{
			name:              "empty string in storeIds - service layer handles validation",
			requestBody:       `{"storeIds": [""]}`,
			expectCode:        http.StatusInternalServerError, // Service layer will handle ID parsing validation
			shouldCallService: true,
		},
		{
			name:              "mixed valid and empty storeIds - service layer handles validation",
			requestBody:       `{"storeIds": ["1", "", "3"]}`,
			expectCode:        http.StatusInternalServerError, // Service layer will handle ID parsing validation
			shouldCallService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For tests that might call the service, set up a mock expectation
			if tt.shouldCallService {
				var req staff.DeleteStoreAccessRequest
				_ = json.Unmarshal([]byte(tt.requestBody), &req)
				mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(nil, errors.New("invalid store IDs"))
			}

			c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", []byte(tt.requestBody))

			handler.DeleteStoreAccess(c)

			assert.Equal(t, tt.expectCode, w.Code)

			var response common.ApiResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Should have message, no data
			assert.NotEmpty(t, response.Message)
			assert.Nil(t, response.Data)
		})
	}

	// Verify service was called for tests that should call it
	mockService.AssertExpectations(t)
}

func TestDeleteStoreAccessHandler_DeleteStoreAccess_ValidRequestWithMultipleStores(t *testing.T) {
	setupTestGinForStoreAccess()

	// Create mock service
	mockService := new(MockDeleteStoreAccessService)
	handler := NewDeleteStoreAccessHandler(mockService)

	// Test data with multiple stores
	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2", "3"},
	}
	reqBody, _ := json.Marshal(req)

	expectedResponse := &staff.DeleteStoreAccessResponse{
		StaffUserID: "987654321",
		StoreList:   []common.Store{}, // No stores left after deletion
	}

	// Set up mock expectations
	mockService.On("DeleteStoreAccess", mock.Anything, "987654321", req, int64(123456789), staff.RoleAdmin, []int64{1, 2}).Return(expectedResponse, nil)

	// Create test context
	c, w := setupTestContextWithStaffForDelete("DELETE", "/api/staff/987654321/store-access", reqBody)

	// Call handler
	handler.DeleteStoreAccess(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Empty(t, response.Message)
	assert.NotNil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Parse the data field
	dataBytes, _ := json.Marshal(response.Data)
	var deleteResponse staff.DeleteStoreAccessResponse
	err = json.Unmarshal(dataBytes, &deleteResponse)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse.StaffUserID, deleteResponse.StaffUserID)
	assert.Equal(t, expectedResponse.StoreList, deleteResponse.StoreList)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}
