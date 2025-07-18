package schedule

import (
	"context"
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
)

type MockDeleteTimeSlotTemplateItemService struct {
	mock.Mock
}

func (m *MockDeleteTimeSlotTemplateItemService) DeleteTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, staffContext common.StaffContext) (*schedule.DeleteTimeSlotTemplateItemResponse, error) {
	args := m.Called(ctx, templateID, itemID, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.DeleteTimeSlotTemplateItemResponse), args.Error(1)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	expectedResponse := &schedule.DeleteTimeSlotTemplateItemResponse{
		Deleted: []string{"6100000003"},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", staffContext).Return(expectedResponse, nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_MissingTemplateID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates//items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: ""},
		{Key: "itemId", Value: "6100000003"},
	}

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_MissingItemID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: ""},
	}

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_MissingStaffContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	// Don't set staffContext

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_PermissionDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleStylist,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.AuthPermissionDenied, "insufficient permissions", nil)
	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", staffContext).Return(nil, serviceError)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid template ID", nil)
	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "invalid", "6100000003", staffContext).Return(nil, serviceError)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/invalid/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "invalid"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_TemplateNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.TimeSlotTemplateNotFound, "time slot template not found", nil)
	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", staffContext).Return(nil, serviceError)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_ItemNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.TimeSlotTemplateItemNotFound, "time slot template item not found", nil)
	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", staffContext).Return(nil, serviceError)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemHandler_DeleteTimeSlotTemplateItem_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDeleteTimeSlotTemplateItemService{}
	handler := NewDeleteTimeSlotTemplateItemHandler(mockService)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.SysDatabaseError, "database error", errors.New("connection failed"))
	mockService.On("DeleteTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", staffContext).Return(nil, serviceError)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/time-slot-templates/6000000011/items/6100000003", nil)

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}
	c.Set("staffContext", staffContext)

	// Call handler
	handler.DeleteTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}