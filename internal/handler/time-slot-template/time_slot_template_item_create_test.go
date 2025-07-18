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

// MockCreateTimeSlotTemplateItemService implements the CreateTimeSlotTemplateItemServiceInterface for testing
type MockCreateTimeSlotTemplateItemService struct {
	mock.Mock
}

// Ensure MockCreateTimeSlotTemplateItemService implements the interface
var _ timeSlotTemplateService.CreateTimeSlotTemplateItemServiceInterface = (*MockCreateTimeSlotTemplateItemService)(nil)

func (m *MockCreateTimeSlotTemplateItemService) CreateTimeSlotTemplateItem(ctx context.Context, templateID string, req timeSlotTemplate.CreateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateItemResponse, error) {
	args := m.Called(ctx, templateID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.CreateTimeSlotTemplateItemResponse), args.Error(1)
}

func setupTestGinForCreateTimeSlotTemplateItem() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestNewCreateTimeSlotTemplateItemHandler(t *testing.T) {
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_Success(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations
	expectedResponse := &timeSlotTemplate.CreateTimeSlotTemplateItemResponse{
		ID:         "6100000003",
		TemplateID: "6000000011",
		StartTime:  "10:00",
		EndTime:    "13:00",
	}

	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(expectedResponse, nil)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

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
	assert.Equal(t, expectedResponse.TemplateID, responseData["templateId"])
	assert.Equal(t, expectedResponse.StartTime, responseData["startTime"])
	assert.Equal(t, expectedResponse.EndTime, responseData["endTime"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_MissingStaffContext(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Call handler
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_MissingTemplateID(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context but missing templateId param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates//items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "templateId為必填項目", response.Errors["templateId"])
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_InvalidJSON(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_MissingRequiredFields(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Create request with missing required fields
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		// Missing startTime and endTime
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_PermissionDenied(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations - return permission denied error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context with insufficient permissions
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
	handler.CreateTimeSlotTemplateItem(c)

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

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_TemplateNotFound(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations - return template not found error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotTemplateNotFound)
	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

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

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_TimeConflict(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations - return time conflict error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleTimeConflict)
	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "時間區段重疊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_ValidationError(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations - return validation error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ValInputValidationFailed)
	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "invalid",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemHandler_CreateTimeSlotTemplateItem_InternalError(t *testing.T) {
	setupTestGinForCreateTimeSlotTemplateItem()

	// Create mock service
	mockService := new(MockCreateTimeSlotTemplateItemService)
	handler := NewCreateTimeSlotTemplateItemHandler(mockService)

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("CreateTimeSlotTemplateItem", mock.Anything, "6000000011", mock.AnythingOfType("timeSlotTemplate.CreateTimeSlotTemplateItemRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, internalError)

	// Create request
	createReq := timeSlotTemplate.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "13:00",
	}
	reqBody, _ := json.Marshal(createReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/time-slot-templates/6000000011/items", bytes.NewBuffer(reqBody))
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
	handler.CreateTimeSlotTemplateItem(c)

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
