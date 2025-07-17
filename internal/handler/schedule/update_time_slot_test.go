package schedule

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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
)

// MockUpdateTimeSlotService implements the UpdateTimeSlotServiceInterface for testing
type MockUpdateTimeSlotService struct {
	mock.Mock
}

// Ensure MockUpdateTimeSlotService implements the interface
var _ scheduleService.UpdateTimeSlotServiceInterface = (*MockUpdateTimeSlotService)(nil)

func (m *MockUpdateTimeSlotService) UpdateTimeSlot(ctx context.Context, scheduleID, timeSlotID string, req schedule.UpdateTimeSlotRequest, staffContext common.StaffContext) (*schedule.UpdateTimeSlotResponse, error) {
	args := m.Called(ctx, scheduleID, timeSlotID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.UpdateTimeSlotResponse), args.Error(1)
}

func setupTestGinForUpdateTimeSlot() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestNewUpdateTimeSlotHandler(t *testing.T) {
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_Success(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations
	isAvailable := false
	expectedResponse := &schedule.UpdateTimeSlotResponse{
		ID:          "5000000001",
		ScheduleID:  "4000000001",
		StartTime:   "09:00",
		EndTime:     "12:00",
		IsAvailable: isAvailable,
	}

	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(expectedResponse, nil)

	// Create request
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

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
	assert.Equal(t, expectedResponse.ScheduleID, responseData["scheduleId"])
	assert.Equal(t, expectedResponse.StartTime, responseData["startTime"])
	assert.Equal(t, expectedResponse.EndTime, responseData["endTime"])
	assert.Equal(t, expectedResponse.IsAvailable, responseData["isAvailable"])

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_MissingStaffContext(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context WITHOUT staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
	}

	// Call handler
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_MissingScheduleID(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context but missing scheduleId param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules//time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "scheduleId為必填項目", response.Errors["scheduleId"])
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_MissingTimeSlotID(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context but missing timeSlotId param
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
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
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
	assert.Equal(t, "timeSlotId為必填項目", response.Errors["timeSlotId"])
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_InvalidJSON(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Create request with invalid JSON
	reqBody := []byte(`{"invalid": json}`)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Errors)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_ValidationError(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return validation error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request with empty fields
	updateReq := schedule.UpdateTimeSlotRequest{
		// All fields are nil
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

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

func TestUpdateTimeSlotHandler_UpdateTimeSlot_TimeSlotNotFound(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return not found error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotNotFound)
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "時段不存在或已被刪除", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_PermissionDenied(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return permission denied error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

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

func TestUpdateTimeSlotHandler_UpdateTimeSlot_TimeSlotAlreadyBooked(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return already booked error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotAlreadyBookedDoNotUpdate)
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "時段已被預約，無法更新", response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotHandler_UpdateTimeSlot_TimeConflict(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return time conflict error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ScheduleTimeConflict)
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, serviceError)

	// Create request
	startTime := "10:00"
	endTime := "14:00"
	updateReq := schedule.UpdateTimeSlotRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

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

func TestUpdateTimeSlotHandler_UpdateTimeSlot_InternalError(t *testing.T) {
	setupTestGinForUpdateTimeSlot()

	// Create mock service
	mockService := new(MockUpdateTimeSlotService)
	handler := NewUpdateTimeSlotHandler(mockService)

	// Set up mock expectations - return internal error
	internalError := errors.New("database connection failed")
	mockService.On("UpdateTimeSlot", mock.Anything, "4000000001", "5000000001", mock.AnythingOfType("schedule.UpdateTimeSlotRequest"), mock.AnythingOfType("common.StaffContext")).Return(nil, internalError)

	// Create request
	isAvailable := false
	updateReq := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}
	reqBody, _ := json.Marshal(updateReq)

	// Create test context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/schedules/4000000001/time-slots/5000000001", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
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
	handler.UpdateTimeSlot(c)

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