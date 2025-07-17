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

// MockCreateStylistService implements the CreateStylistServiceInterface for testing
type MockCreateStylistService struct {
	mock.Mock
}

// Ensure MockCreateStylistService implements the interface
var _ stylistService.CreateStylistServiceInterface = (*MockCreateStylistService)(nil)

func (m *MockCreateStylistService) CreateStylist(ctx context.Context, req stylist.CreateStylistRequest, staffUserID int64) (*stylist.CreateStylistResponse, error) {
	args := m.Called(ctx, req, staffUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stylist.CreateStylistResponse), args.Error(1)
}

func setupTestGinForStylist() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestCreateStylistHandler_CreateStylist_Success(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockCreateStylistService)
	handler := NewCreateStylistHandler(mockService)

	// Create test request
	req := stylist.CreateStylistRequest{
		StylistName:  "Jane 美甲師",
		GoodAtShapes: []string{"方形", "圓形"},
		GoodAtColors: []string{"裸色系", "粉嫩系"},
		GoodAtStyles: []string{"手繪", "簡約"},
		IsIntrovert:  boolPtr(false),
	}

	// Create expected response
	expectedResponse := &stylist.CreateStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Jane 美甲師",
		GoodAtShapes: []string{"方形", "圓形"},
		GoodAtColors: []string{"裸色系", "粉嫩系"},
		GoodAtStyles: []string{"手繪", "簡約"},
		IsIntrovert:  false,
	}

	// Mock service call
	mockService.On("CreateStylist", mock.Anything, req, int64(12345)).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/stylists", bytes.NewBuffer(jsonData))
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
	handler.CreateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

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
	assert.Equal(t, "Jane 美甲師", responseData["stylistName"])

	mockService.AssertExpectations(t)
}

func TestCreateStylistHandler_CreateStylist_ValidationError(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockCreateStylistService)
	handler := NewCreateStylistHandler(mockService)

	// Create invalid request (missing required name)
	req := stylist.CreateStylistRequest{
		// Name is required but missing
		GoodAtShapes: []string{"方形"},
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/stylists", bytes.NewBuffer(jsonData))
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
	handler.CreateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response["message"])
	assert.Contains(t, response, "errors")

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateStylist")
}

func TestCreateStylistHandler_CreateStylist_InvalidJSON(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockCreateStylistService)
	handler := NewCreateStylistHandler(mockService)

	// Create invalid JSON
	invalidJSON := `{"name": "Jane 美甲師", "goodAtShapes": [`

	// Create HTTP request
	httpReq := httptest.NewRequest(http.MethodPost, "/api/stylists", bytes.NewBufferString(invalidJSON))
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
	handler.CreateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "輸入驗證失敗", response["message"])

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateStylist")
}

func TestCreateStylistHandler_CreateStylist_NoStaffContext(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockCreateStylistService)
	handler := NewCreateStylistHandler(mockService)

	// Create test request
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/stylists", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Setup Gin context WITHOUT staff context
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call handler
	handler.CreateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "未找到使用者認證資訊", response["message"])

	// Service should not be called
	mockService.AssertNotCalled(t, "CreateStylist")
}

func TestCreateStylistHandler_CreateStylist_ServiceError(t *testing.T) {
	setupTestGinForStylist()

	// Create mock service
	mockService := new(MockCreateStylistService)
	handler := NewCreateStylistHandler(mockService)

	// Create test request
	req := stylist.CreateStylistRequest{
		StylistName: "Jane 美甲師",
	}

	// Mock service error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.StylistAlreadyExists)
	mockService.On("CreateStylist", mock.Anything, req, int64(12345)).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/stylists", bytes.NewBuffer(jsonData))
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
	handler.CreateStylist(c)

	// Assert response
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "該員工已建立過美甲師資料，請使用修改功能", response["message"])

	mockService.AssertExpectations(t)
}

// Helper function to create boolean pointer
func boolPtr(b bool) *bool {
	return &b
}
