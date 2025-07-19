package store

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

func init() {
	// Initialize error manager for testing
	manager := errorCodes.GetManager()
	wd, _ := os.Getwd()
	// Navigate up from handler/store to project root
	errorFilePath := filepath.Join(wd, "..", "..", "..", "internal", "errors", "errors.yaml")
	if err := manager.LoadFromFile(errorFilePath); err != nil {
		panic("Failed to load errors.yaml for testing: " + err.Error())
	}

	// register custom validators for testing
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("taiwanlandline", utils.ValidateTaiwanLandline)
	}
}

func TestNewUpdateStoreHandler(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockUpdateStoreService is a mock implementation of UpdateStoreServiceInterface
type MockUpdateStoreService struct {
	mock.Mock
}

func (m *MockUpdateStoreService) UpdateStore(ctx context.Context, storeID string, req store.UpdateStoreRequest, staffContext common.StaffContext) (*store.UpdateStoreResponse, error) {
	args := m.Called(ctx, storeID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.UpdateStoreResponse), args.Error(1)
}

func TestUpdateStoreHandler_UpdateStore_Success(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	updateReq := store.UpdateStoreRequest{
		Name:    stringPtr("Updated Store"),
		Address: stringPtr("Updated Address"),
		Phone:   stringPtr("02-87654321"),
	}
	reqBody, _ := json.Marshal(updateReq)

	c.Request, _ = http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Updated Store",
		Address:  "Updated Address",
		Phone:    "02-87654321",
		IsActive: true,
	}
	mockService.On("UpdateStore", mock.Anything, "8000000001", mock.AnythingOfType("store.UpdateStoreRequest"), staffContext).Return(expectedResponse, nil)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_InvalidJSON(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with invalid JSON
	reqBody := `{"name":}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_ValidationError(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with validation error (name too long)
	reqBody := `{"name":"` + stringRepeat("a", 101) + `"}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_InvalidTaiwanLandline(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with invalid Taiwan landline
	reqBody := `{"phone":"invalid-phone"}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_MissingStoreId(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without store ID param
	reqBody := `{"name":"Updated Store"}`
	req, _ := http.NewRequest("PATCH", "/api/stores/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	// No params set - missing storeId

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_NoUpdates(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with no fields to update
	reqBody := `{}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_MissingStaffContext(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without staff context
	reqBody := `{"name":"Updated Store"}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_ServiceError(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	reqBody := `{"name":"Duplicate Name"}`
	req, _ := http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service error
	serviceError := errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "store name already exists", nil)
	mockService.On("UpdateStore", mock.Anything, "8000000001", mock.AnythingOfType("store.UpdateStoreRequest"), staffContext).Return(nil, serviceError)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func stringRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func TestUpdateStoreHandler_UpdateStore_Success_UpdateIsActive(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request to deactivate store
	updateReq := store.UpdateStoreRequest{
		IsActive: boolPtr(false),
	}
	reqBody, _ := json.Marshal(updateReq)

	c.Request, _ = http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Test Store",
		Address:  "Test Address",
		Phone:    "02-12345678",
		IsActive: false,
	}
	mockService.On("UpdateStore", mock.Anything, "8000000001", mock.AnythingOfType("store.UpdateStoreRequest"), staffContext).Return(expectedResponse, nil)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response to verify IsActive was updated
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	data := response["data"].(map[string]interface{})
	assert.False(t, data["isActive"].(bool))
	
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_Success_ReactivateStore(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request to reactivate store
	updateReq := store.UpdateStoreRequest{
		IsActive: boolPtr(true),
	}
	reqBody, _ := json.Marshal(updateReq)

	c.Request, _ = http.NewRequest("PATCH", "/api/stores/8000000002", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000002"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000002",
		Name:     "Reactivated Store",
		Address:  "Test Address",
		Phone:    "02-87654321",
		IsActive: true,
	}
	mockService.On("UpdateStore", mock.Anything, "8000000002", mock.AnythingOfType("store.UpdateStoreRequest"), staffContext).Return(expectedResponse, nil)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response to verify IsActive was updated
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	data := response["data"].(map[string]interface{})
	assert.True(t, data["isActive"].(bool))
	
	mockService.AssertExpectations(t)
}

func TestUpdateStoreHandler_UpdateStore_Success_MultipleFieldsWithIsActive(t *testing.T) {
	mockService := &MockUpdateStoreService{}
	handler := NewUpdateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with multiple fields including IsActive
	updateReq := store.UpdateStoreRequest{
		Name:     stringPtr("Updated Store"),
		Address:  stringPtr("Updated Address"),
		Phone:    stringPtr("02-98765432"),
		IsActive: boolPtr(false),
	}
	reqBody, _ := json.Marshal(updateReq)

	c.Request, _ = http.NewRequest("PATCH", "/api/stores/8000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "storeId", Value: "8000000001"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Updated Store",
		Address:  "Updated Address",
		Phone:    "02-98765432",
		IsActive: false,
	}
	mockService.On("UpdateStore", mock.Anything, "8000000001", mock.AnythingOfType("store.UpdateStoreRequest"), staffContext).Return(expectedResponse, nil)

	handler.UpdateStore(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response to verify all fields were updated
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Updated Store", data["name"].(string))
	assert.Equal(t, "Updated Address", data["address"].(string))
	assert.Equal(t, "02-98765432", data["phone"].(string))
	assert.False(t, data["isActive"].(bool))
	
	mockService.AssertExpectations(t)
}