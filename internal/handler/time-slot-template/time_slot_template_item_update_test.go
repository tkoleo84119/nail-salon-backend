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
)

type MockUpdateTimeSlotTemplateItemService struct {
	mock.Mock
}

func (m *MockUpdateTimeSlotTemplateItemService) UpdateTimeSlotTemplateItem(ctx context.Context, templateID string, itemID string, req timeSlotTemplate.UpdateTimeSlotTemplateItemRequest, staffContext common.StaffContext) (*timeSlotTemplate.UpdateTimeSlotTemplateItemResponse, error) {
	args := m.Called(ctx, templateID, itemID, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.UpdateTimeSlotTemplateItemResponse), args.Error(1)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateItemResponse{
		ID:         "6100000003",
		TemplateID: "6000000011",
		StartTime:  "14:00",
		EndTime:    "18:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(expectedResponse, nil)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "6100000003", data["id"])
	assert.Equal(t, "6000000011", data["templateId"])
	assert.Equal(t, "14:00", data["startTime"])
	assert.Equal(t, "18:00", data["endTime"])

	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_MissingTemplateID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates//items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: ""},
		{Key: "itemId", Value: "6100000003"},
	}

	// Call handler
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_MissingItemID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: ""},
	}

	// Call handler
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	// Create request with invalid JSON
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "templateId", Value: "6000000011"},
		{Key: "itemId", Value: "6100000003"},
	}

	// Call handler
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_MissingStaffContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_PermissionDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleStylist,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.AuthPermissionDenied, "insufficient permissions", nil)
	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(nil, serviceError)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusForbidden, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "invalid",
		EndTime:   "18:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid start time format", nil)
	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(nil, serviceError)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_ItemNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.TimeSlotTemplateItemNotFound, "time slot template item not found", nil)
	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(nil, serviceError)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_TimeConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "14:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.ScheduleTimeConflict, "time slot overlaps with existing template item", nil)
	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(nil, serviceError)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemHandler_UpdateTimeSlotTemplateItem_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUpdateTimeSlotTemplateItemService{}
	handler := NewUpdateTimeSlotTemplateItemHandler(mockService)

	request := timeSlotTemplate.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	serviceError := errorCodes.NewServiceError(errorCodes.SysDatabaseError, "database error", errors.New("connection failed"))
	mockService.On("UpdateTimeSlotTemplateItem", mock.Anything, "6000000011", "6100000003", request, staffContext).Return(nil, serviceError)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PATCH", "/api/time-slot-templates/6000000011/items/6100000003", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

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
	handler.UpdateTimeSlotTemplateItem(c)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
