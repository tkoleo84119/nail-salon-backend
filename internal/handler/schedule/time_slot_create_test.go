package schedule

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
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
)

// MockCreateTimeSlotService is a mock implementation of CreateTimeSlotServiceInterface
type MockCreateTimeSlotService struct {
	mock.Mock
}

func (m *MockCreateTimeSlotService) CreateTimeSlot(ctx context.Context, scheduleID string, req scheduleModel.CreateTimeSlotRequest, staffContext common.StaffContext) (*scheduleModel.CreateTimeSlotResponse, error) {
	args := m.Called(ctx, scheduleID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*scheduleModel.CreateTimeSlotResponse), args.Error(1)
}

func TestCreateTimeSlotHandler_CreateTimeSlot_Success(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	// Expected response
	expectedResponse := scheduleModel.CreateTimeSlotResponse{
		ID:         "5000000001",
		ScheduleID: "4000000001",
		StartTime:  "09:00",
		EndTime:    "12:00",
	}

	// Setup mock
	mockService.On("CreateTimeSlot", mock.Anything, "4000000001", request, staffContext).Return(&expectedResponse, nil)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotHandler_CreateTimeSlot_NoStaffContext(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request without staff context
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Setup Gin context without staff context
	c, w := setupTestContext(nil)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 401 for missing staff context
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateTimeSlot")
}

func TestCreateTimeSlotHandler_CreateTimeSlot_InvalidJSON(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup Gin context with invalid JSON
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 400 for invalid JSON
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateTimeSlot")
}

func TestCreateTimeSlotHandler_CreateTimeSlot_MissingScheduleID(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup request
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Setup Gin context without schedule ID param
	c, w := setupTestContext(&staffContext)
	// No schedule ID param set

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules//time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 400 for validation error
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateTimeSlot")
}

func TestCreateTimeSlotHandler_CreateTimeSlot_ValidationError(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request with validation errors (missing required fields)
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "", // Missing start time
		EndTime:   "", // Missing end time
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 400 for validation error
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateTimeSlot")
}

func TestCreateTimeSlotHandler_CreateTimeSlot_ServiceError(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup mock to return error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleNotFound)
	mockService.On("CreateTimeSlot", mock.Anything, "4000000001", request, staffContext).Return(nil, serviceError)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 404 for SCHEDULE_NOT_FOUND error
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotHandler_CreateTimeSlot_TimeSlotOverlap(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup mock to return overlap error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleTimeConflict)
	mockService.On("CreateTimeSlot", mock.Anything, "4000000001", request, staffContext).Return(nil, serviceError)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 400 for SCHEDULE_TIME_CONFLICT error
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotHandler_CreateTimeSlot_PermissionDenied(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	// Setup request
	request := scheduleModel.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleStylist, // Stylist trying to modify other's schedule
	}

	// Setup mock to return permission denied
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("CreateTimeSlot", mock.Anything, "4000000001", request, staffContext).Return(nil, serviceError)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	c.Params = []gin.Param{{Key: "scheduleId", Value: "4000000001"}}

	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/4000000001/time-slots", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateTimeSlot(c)

	// Assert - expect 403 for AUTH_PERMISSION_DENIED error
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	mockService.AssertExpectations(t)
}

func TestNewCreateTimeSlotHandler(t *testing.T) {
	mockService := &MockCreateTimeSlotService{}
	handler := NewCreateTimeSlotHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}