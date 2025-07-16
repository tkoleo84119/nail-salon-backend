package stylist

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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
)

// MockUpdateStylistService implements the UpdateStylistServiceInterface for testing
type MockUpdateStylistService struct {
	mock.Mock
}

// Ensure MockUpdateStylistService implements the interface
var _ stylistService.UpdateStylistServiceInterface = (*MockUpdateStylistService)(nil)

func (m *MockUpdateStylistService) UpdateStylist(ctx context.Context, req stylist.UpdateStylistRequest, staffUserID int64) (*stylist.UpdateStylistResponse, error) {
	args := m.Called(ctx, req, staffUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stylist.UpdateStylistResponse), args.Error(1)
}

func TestUpdateStylistHandler_UpdateStylist_Success(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create test request
	stylistName := "Jane Updated"
	goodAtShapes := []string{"橢圓形", "方形"}
	isIntrovert := true

	req := stylist.UpdateStylistRequest{
		StylistName:  &stylistName,
		GoodAtShapes: &goodAtShapes,
		IsIntrovert:  &isIntrovert,
	}

	// Create expected response
	expectedResponse := &stylist.UpdateStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Jane Updated",
		GoodAtShapes: []string{"橢圓形", "方形"},
		GoodAtColors: []string{"裸色系"},
		GoodAtStyles: []string{"簡約"},
		IsIntrovert:  true,
	}

	// Mock service call
	mockService.On("UpdateStylist", mock.Anything, req, int64(12345)).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context with staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "testuser",
		Role:     staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.NotNil(t, response.Data)

	// Verify response data
	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "18000000001", responseData["id"])
	assert.Equal(t, "12345", responseData["staffUserId"])
	assert.Equal(t, "Jane Updated", responseData["stylistName"])
	assert.Equal(t, true, responseData["isIntrovert"])

	mockService.AssertExpectations(t)
}

func TestUpdateStylistHandler_UpdateStylist_ValidationError(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create invalid request (stylistName too long)
	longName := string(make([]byte, 51)) // 51 characters, exceeds max of 50
	req := stylist.UpdateStylistRequest{
		StylistName: &longName,
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context with staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "testuser",
		Role:     staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response["message"])
	assert.Contains(t, response, "errors")

	// Service should not be called
	mockService.AssertNotCalled(t, "UpdateStylist")
}

func TestUpdateStylistHandler_UpdateStylist_InvalidJSON(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create invalid JSON
	invalidJSON := `{"stylistName": "Jane Updated", "goodAtShapes": [`

	// Create HTTP request
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBufferString(invalidJSON))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context with staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "testuser",
		Role:     staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "JSON格式錯誤", response["message"])

	// Service should not be called
	mockService.AssertNotCalled(t, "UpdateStylist")
}

func TestUpdateStylistHandler_UpdateStylist_NoStaffContext(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create test request
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context WITHOUT staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response["message"])

	// Service should not be called
	mockService.AssertNotCalled(t, "UpdateStylist")
}

func TestUpdateStylistHandler_UpdateStylist_ServiceError(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create test request
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock service error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotCreated)
	mockService.On("UpdateStylist", mock.Anything, req, int64(12345)).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context with staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "testuser",
		Role:     staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "尚未建立美甲師資料，請先新增", response["message"])

	mockService.AssertExpectations(t)
}

func TestUpdateStylistHandler_UpdateStylist_EmptyRequest(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockUpdateStylistService)
	handler := NewUpdateStylistHandler(mockService)

	// Create empty request
	req := stylist.UpdateStylistRequest{}

	// Mock service error for empty fields
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	mockService.On("UpdateStylist", mock.Anything, req, int64(12345)).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPatch, "/api/stylists/me", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context with staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "testuser",
		Role:     staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Call handler
	handler.UpdateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "至少需要提供一個欄位進行更新", response["message"])

	mockService.AssertExpectations(t)
}