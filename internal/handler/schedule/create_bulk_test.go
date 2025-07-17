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
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
)

// MockCreateSchedulesBulkService is a mock implementation of CreateSchedulesBulkServiceInterface
type MockCreateSchedulesBulkService struct {
	mock.Mock
}

func (m *MockCreateSchedulesBulkService) CreateSchedulesBulk(ctx context.Context, req schedule.CreateSchedulesBulkRequest, staffContext common.StaffContext) (*schedule.CreateSchedulesBulkResponse, error) {
	args := m.Called(ctx, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.CreateSchedulesBulkResponse), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func setupTestContext(staffContext *common.StaffContext) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	if staffContext != nil {
		c.Set(middleware.UserContextKey, *staffContext)
	}
	
	return c, w
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_Success(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup request
	note := "Test note"
	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{
			{
				WorkDate: "2023-12-01",
				Note:     &note,
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "09:00", EndTime: "10:00"},
					{StartTime: "14:00", EndTime: "15:00"},
				},
			},
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	// Expected response
	expectedResponse := schedule.CreateSchedulesBulkResponse{
		{
			ScheduleID: "99999",
			StylistID:  "12345",
			StoreID:    "67890",
			WorkDate:   "2023-12-01",
			Note:       &note,
			TimeSlots: []schedule.TimeSlotResponse{
				{ID: "11111", StartTime: "09:00", EndTime: "10:00"},
				{ID: "22222", StartTime: "14:00", EndTime: "15:00"},
			},
		},
	}

	// Setup mock
	mockService.On("CreateSchedulesBulk", mock.Anything, request, staffContext).Return(&expectedResponse, nil)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_NoStaffContext(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup request without staff context
	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	// Setup Gin context without staff context
	c, w := setupTestContext(nil)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert - expect 401 for missing staff context
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateSchedulesBulk")
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_InvalidJSON(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup Gin context with invalid JSON
	c, w := setupTestContext(&staffContext)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert - expect 400 for invalid JSON
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateSchedulesBulk")
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_ValidationError(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup request with validation errors (missing required fields)
	request := schedule.CreateSchedulesBulkRequest{
		// StylistID missing
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{}, // Empty schedules (violates min=1)
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert - expect 400 for validation errors
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateSchedulesBulk")
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_ServiceError(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup request
	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{
			{
				WorkDate: "2023-12-01",
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "09:00", EndTime: "10:00"},
				},
			},
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Setup mock to return error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
	mockService.On("CreateSchedulesBulk", mock.Anything, request, staffContext).Return(nil, serviceError)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code) // expect 404 for StylistNotFound
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	mockService.AssertExpectations(t)
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_PermissionDenied(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup request
	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{
			{
				WorkDate: "2023-12-01",
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "09:00", EndTime: "10:00"},
				},
			},
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleStylist, // Stylist trying to create for different stylist
	}

	// Setup mock to return permission denied
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("CreateSchedulesBulk", mock.Anything, request, staffContext).Return(nil, serviceError)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code) // expect 403 for permission denied
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Message)

	mockService.AssertExpectations(t)
}

func TestCreateSchedulesBulkHandler_CreateSchedulesBulk_ComplexRequest(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	// Setup complex request with multiple schedules
	note1 := "Morning session"
	note2 := "Afternoon session"
	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{
			{
				WorkDate: "2023-12-01",
				Note:     &note1,
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "09:00", EndTime: "10:00"},
					{StartTime: "10:30", EndTime: "11:30"},
				},
			},
			{
				WorkDate: "2023-12-02",
				Note:     &note2,
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "14:00", EndTime: "15:00"},
				},
			},
			{
				WorkDate: "2023-12-03",
				Note:     nil, // No note
				TimeSlots: []schedule.TimeSlotRequest{
					{StartTime: "16:00", EndTime: "17:00"},
					{StartTime: "17:30", EndTime: "18:30"},
				},
			},
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleManager,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	// Expected response
	expectedResponse := schedule.CreateSchedulesBulkResponse{
		{
			ScheduleID: "99991",
			StylistID:  "12345",
			StoreID:    "67890",
			WorkDate:   "2023-12-01",
			Note:       &note1,
			TimeSlots: []schedule.TimeSlotResponse{
				{ID: "11111", StartTime: "09:00", EndTime: "10:00"},
				{ID: "11112", StartTime: "10:30", EndTime: "11:30"},
			},
		},
		{
			ScheduleID: "99992",
			StylistID:  "12345",
			StoreID:    "67890",
			WorkDate:   "2023-12-02",
			Note:       &note2,
			TimeSlots: []schedule.TimeSlotResponse{
				{ID: "11113", StartTime: "14:00", EndTime: "15:00"},
			},
		},
		{
			ScheduleID: "99993",
			StylistID:  "12345",
			StoreID:    "67890",
			WorkDate:   "2023-12-03",
			Note:       nil,
			TimeSlots: []schedule.TimeSlotResponse{
				{ID: "11114", StartTime: "16:00", EndTime: "17:00"},
				{ID: "11115", StartTime: "17:30", EndTime: "18:30"},
			},
		},
	}

	// Setup mock
	mockService.On("CreateSchedulesBulk", mock.Anything, request, staffContext).Return(&expectedResponse, nil)

	// Setup Gin context
	c, w := setupTestContext(&staffContext)
	
	reqJSON, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/api/schedules/bulk", bytes.NewBuffer(reqJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateSchedulesBulk(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Verify response data structure
	responseData, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, responseData, 3) // 3 schedules

	mockService.AssertExpectations(t)
}

func TestNewCreateSchedulesBulkHandler(t *testing.T) {
	mockService := &MockCreateSchedulesBulkService{}
	handler := NewCreateSchedulesBulkHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}