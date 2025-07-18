package timeSlotTemplate

import (
	"bytes"
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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
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

func TestNewCreateTimeSlotTemplateHandler(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := timeSlotTemplateService.NewCreateTimeSlotTemplateService(mockQuerier, nil)
	handler := NewCreateTimeSlotTemplateHandler(service)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockCreateTimeSlotTemplateService is a mock implementation of CreateTimeSlotTemplateServiceInterface
type MockCreateTimeSlotTemplateService struct {
	mock.Mock
}

func (m *MockCreateTimeSlotTemplateService) CreateTimeSlotTemplate(ctx context.Context, req timeSlotTemplate.CreateTimeSlotTemplateRequest, staffContext common.StaffContext) (*timeSlotTemplate.CreateTimeSlotTemplateResponse, error) {
	args := m.Called(ctx, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.CreateTimeSlotTemplateResponse), args.Error(1)
}

func TestCreateTimeSlotTemplateHandler_CreateTimeSlotTemplate_MissingStaffContext(t *testing.T) {
	mockService := &MockCreateTimeSlotTemplateService{}
	handler := NewCreateTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without staff context
	reqBody := `{"name":"Test Template","timeSlots":[{"startTime":"09:00","endTime":"12:00"}]}`
	req, _ := http.NewRequest("POST", "/api/time-slot-templates", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	handler.CreateTimeSlotTemplate(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateTimeSlotTemplateHandler_CreateTimeSlotTemplate_InvalidJSON(t *testing.T) {
	mockService := &MockCreateTimeSlotTemplateService{}
	handler := NewCreateTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with staff context but invalid JSON
	reqBody := `{"name":"Test Template","timeSlots":[{"startTime":"09:00","endTime":}]}` // Invalid JSON
	req, _ := http.NewRequest("POST", "/api/time-slot-templates", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.CreateTimeSlotTemplate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTimeSlotTemplateHandler_CreateTimeSlotTemplate_ValidationError(t *testing.T) {
	mockService := &MockCreateTimeSlotTemplateService{}
	handler := NewCreateTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with validation errors (missing required fields)
	reqBody := `{"timeSlots":[{"startTime":"09:00","endTime":"12:00"}]}` // Missing name
	req, _ := http.NewRequest("POST", "/api/time-slot-templates", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.CreateTimeSlotTemplate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTimeSlotTemplateHandler_CreateTimeSlotTemplate_ServiceError(t *testing.T) {
	mockService := &MockCreateTimeSlotTemplateService{}
	handler := NewCreateTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	reqBody := `{"name":"Test Template","timeSlots":[{"startTime":"09:00","endTime":"12:00"}]}`
	req, _ := http.NewRequest("POST", "/api/time-slot-templates", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service to return error
	expectedRequest := timeSlotTemplate.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		TimeSlots: []timeSlotTemplate.TimeSlotItem{
			{StartTime: "09:00", EndTime: "12:00"},
		},
	}
	serviceErr := errorCodes.NewServiceError(errorCodes.AuthPermissionDenied, "permission denied", nil)
	mockService.On("CreateTimeSlotTemplate", mock.Anything, expectedRequest, staffContext).Return(nil, serviceErr)

	handler.CreateTimeSlotTemplate(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateHandler_CreateTimeSlotTemplate_Success(t *testing.T) {
	mockService := &MockCreateTimeSlotTemplateService{}
	handler := NewCreateTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	reqBody := `{"name":"Test Template","note":"Test note","timeSlots":[{"startTime":"09:00","endTime":"12:00"},{"startTime":"13:00","endTime":"17:00"}]}`
	req, _ := http.NewRequest("POST", "/api/time-slot-templates", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service to return success
	expectedRequest := timeSlotTemplate.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		Note: "Test note",
		TimeSlots: []timeSlotTemplate.TimeSlotItem{
			{StartTime: "09:00", EndTime: "12:00"},
			{StartTime: "13:00", EndTime: "17:00"},
		},
	}
	expectedResponse := &timeSlotTemplate.CreateTimeSlotTemplateResponse{
		ID:   "6000000001",
		Name: "Test Template",
		Note: "Test note",
		TimeSlots: []timeSlotTemplate.TimeSlotItemResponse{
			{ID: "6100000001", StartTime: "09:00", EndTime: "12:00"},
			{ID: "6100000002", StartTime: "13:00", EndTime: "17:00"},
		},
	}
	mockService.On("CreateTimeSlotTemplate", mock.Anything, expectedRequest, staffContext).Return(expectedResponse, nil)

	handler.CreateTimeSlotTemplate(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}
