package schedule

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func init() {
	// Initialize error manager for testing
	manager := errorCodes.GetManager()
	wd, _ := os.Getwd()
	// Navigate up from handler/schedule to project root
	errorFilePath := filepath.Join(wd, "..", "..", "..", "internal", "errors", "errors.yaml")
	if err := manager.LoadFromFile(errorFilePath); err != nil {
		panic("Failed to load errors.yaml for testing: " + err.Error())
	}
}

func TestNewDeleteTimeSlotHandler(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := scheduleService.NewDeleteTimeSlotService(mockQuerier)
	handler := NewDeleteTimeSlotHandler(service)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockDeleteTimeSlotService is a mock implementation of DeleteTimeSlotServiceInterface
type MockDeleteTimeSlotService struct {
	mock.Mock
}

func (m *MockDeleteTimeSlotService) DeleteTimeSlot(ctx context.Context, scheduleID string, timeSlotID string, staffContext common.StaffContext) (*schedule.DeleteTimeSlotResponse, error) {
	args := m.Called(ctx, scheduleID, timeSlotID, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.DeleteTimeSlotResponse), args.Error(1)
}

func TestDeleteTimeSlotHandler_DeleteTimeSlot_MissingStaffContext(t *testing.T) {
	mockService := &MockDeleteTimeSlotService{}
	handler := NewDeleteTimeSlotHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without staff context
	req, _ := http.NewRequest("DELETE", "/api/schedules/4000000001/time-slots/5000000001", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
	}

	handler.DeleteTimeSlot(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteTimeSlotHandler_DeleteTimeSlot_MissingScheduleID(t *testing.T) {
	mockService := &MockDeleteTimeSlotService{}
	handler := NewDeleteTimeSlotHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with staff context but missing scheduleId (empty string)
	req, _ := http.NewRequest("DELETE", "/api/schedules//time-slots/5000000001", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "scheduleId", Value: ""},
		{Key: "timeSlotId", Value: "5000000001"},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.DeleteTimeSlot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTimeSlotHandler_DeleteTimeSlot_MissingTimeSlotID(t *testing.T) {
	mockService := &MockDeleteTimeSlotService{}
	handler := NewDeleteTimeSlotHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with staff context but missing timeSlotId (empty string)
	req, _ := http.NewRequest("DELETE", "/api/schedules/4000000001/time-slots/", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: ""},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.DeleteTimeSlot(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTimeSlotHandler_DeleteTimeSlot_ServiceError(t *testing.T) {
	mockService := &MockDeleteTimeSlotService{}
	handler := NewDeleteTimeSlotHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	req, _ := http.NewRequest("DELETE", "/api/schedules/4000000001/time-slots/5000000001", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service to return error
	serviceErr := errorCodes.NewServiceError(errorCodes.TimeSlotNotFound, "time slot not found", nil)
	mockService.On("DeleteTimeSlot", mock.Anything, "4000000001", "5000000001", staffContext).Return(nil, serviceErr)

	handler.DeleteTimeSlot(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotHandler_DeleteTimeSlot_Success(t *testing.T) {
	mockService := &MockDeleteTimeSlotService{}
	handler := NewDeleteTimeSlotHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	req, _ := http.NewRequest("DELETE", "/api/schedules/4000000001/time-slots/5000000001", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "scheduleId", Value: "4000000001"},
		{Key: "timeSlotId", Value: "5000000001"},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	// Mock service to return success
	expectedResponse := &schedule.DeleteTimeSlotResponse{
		Deleted: []string{"5000000001"},
	}
	mockService.On("DeleteTimeSlot", mock.Anything, "4000000001", "5000000001", staffContext).Return(expectedResponse, nil)

	handler.DeleteTimeSlot(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}