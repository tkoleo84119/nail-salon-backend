package timeSlotTemplate

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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func init() {
	// Initialize error manager for testing
	manager := errorCodes.GetManager()
	wd, _ := os.Getwd()
	// Navigate up from handler/time-slot-template to project root
	errorFilePath := filepath.Join(wd, "..", "..", "..", "internal", "errors", "errors.yaml")
	if err := manager.LoadFromFile(errorFilePath); err != nil {
		panic("Failed to load errors.yaml for testing: " + err.Error())
	}
}

func TestNewDeleteTimeSlotTemplateHandler(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := timeSlotTemplateService.NewDeleteTimeSlotTemplateService(mockQuerier)
	handler := NewDeleteTimeSlotTemplateHandler(service)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockDeleteTimeSlotTemplateService is a mock implementation of DeleteTimeSlotTemplateServiceInterface
type MockDeleteTimeSlotTemplateService struct {
	mock.Mock
}

func (m *MockDeleteTimeSlotTemplateService) DeleteTimeSlotTemplate(ctx context.Context, templateID string, staffContext common.StaffContext) (*timeSlotTemplate.DeleteTimeSlotTemplateResponse, error) {
	args := m.Called(ctx, templateID, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.DeleteTimeSlotTemplateResponse), args.Error(1)
}

func TestDeleteTimeSlotTemplateHandler_DeleteTimeSlotTemplate_Success(t *testing.T) {
	mockService := &MockDeleteTimeSlotTemplateService{}
	handler := NewDeleteTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &timeSlotTemplate.DeleteTimeSlotTemplateResponse{
		Deleted: []string{"6000000011"},
	}
	mockService.On("DeleteTimeSlotTemplate", mock.Anything, "6000000011", staffContext).Return(expectedResponse, nil)

	handler.DeleteTimeSlotTemplate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateHandler_DeleteTimeSlotTemplate_MissingTemplateID(t *testing.T) {
	mockService := &MockDeleteTimeSlotTemplateService{}
	handler := NewDeleteTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without templateId parameter
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/", nil)
	c.Request = req
	// No templateId parameter set

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.DeleteTimeSlotTemplate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateHandler_DeleteTimeSlotTemplate_MissingStaffContext(t *testing.T) {
	mockService := &MockDeleteTimeSlotTemplateService{}
	handler := NewDeleteTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// No staff context set

	handler.DeleteTimeSlotTemplate(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateHandler_DeleteTimeSlotTemplate_ServiceError(t *testing.T) {
	mockService := &MockDeleteTimeSlotTemplateService{}
	handler := NewDeleteTimeSlotTemplateHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011", nil)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
	}

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service error
	serviceError := errorCodes.NewServiceError(errorCodes.TimeSlotTemplateNotFound, "template not found", nil)
	mockService.On("DeleteTimeSlotTemplate", mock.Anything, "6000000011", staffContext).Return(nil, serviceError)

	handler.DeleteTimeSlotTemplate(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}