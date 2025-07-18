package timeSlotTemplate

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
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
)

// MockUpdateTimeSlotTemplateService implements the UpdateTimeSlotTemplateServiceInterface for testing
type MockUpdateTimeSlotTemplateService struct {
	mock.Mock
}

// Ensure MockUpdateTimeSlotTemplateService implements the interface
var _ timeSlotTemplateService.UpdateTimeSlotTemplateServiceInterface = (*MockUpdateTimeSlotTemplateService)(nil)

func (m *MockUpdateTimeSlotTemplateService) UpdateTimeSlotTemplate(ctx context.Context, templateID string, req timeSlotTemplate.UpdateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateResponse, error) {
	args := m.Called(ctx, templateID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.UpdateTimeSlotTemplateResponse), args.Error(1)
}

func setupTestGinForUpdateTimeSlotTemplate() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestNewUpdateTimeSlotTemplateHandler(t *testing.T) {
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_Success(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations
	name := "Updated Template"
	note := "Updated note"
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateResponse{
		ID:   "6000000011",
		Name: "Updated Template",
		Note: "Updated note",
	}

	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(expectedResponse, nil)

	// Create request
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
		Note: &note,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

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
	assert.Equal(t, expectedResponse.Name, responseData["name"])
	assert.Equal(t, expectedResponse.Note, responseData["note"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_UpdateNameOnly(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations
	name := "Updated Template"
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateResponse{
		ID:   "6000000011",
		Name: "Updated Template",
		Note: "Original note",
	}

	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(expectedResponse, nil)

	// Create request
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

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
	assert.Equal(t, expectedResponse.Name, responseData["name"])
	assert.Equal(t, expectedResponse.Note, responseData["note"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_MissingStaffContext(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Create request
	name := "Updated Template"
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_MissingTemplateID(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Create request
	name := "Updated Template"
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context but missing templateId param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		// Missing templateId param
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "templateId為必填項目", response.Errors[("templateId")])
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_InvalidJSON(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_ValidationError(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations - return validation error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request with empty fields
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		// All fields are nil
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "至少需要提供一個欄位進行更新", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_TemplateNotFound(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations - return not found error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	name := "Updated Template"
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "範本不存在或已被刪除", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_PermissionDenied(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations - return permission denied error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	name := "Updated Template"
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "stylist",
		Role:     staff.RoleStylist,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

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

func TestUpdateTimeSlotTemplateHandler_UpdateTimeSlotTemplate_InternalError(t *testing.T) {
	setupTestGinForUpdateTimeSlotTemplate()

	// Create mock service
	mockService := new(MockUpdateTimeSlotTemplateService)
	handler := NewUpdateTimeSlotTemplateHandler(mockService)

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("UpdateTimeSlotTemplate", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.UpdateTimeSlotTemplateRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, internalError)

	// Create request
	name := "Updated Template"
	updateReq := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/api/time-slot-templates/6000000011", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "11111",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateTimeSlotTemplate(c)

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
