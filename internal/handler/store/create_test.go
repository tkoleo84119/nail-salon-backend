package store

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

func init() {
	// Initialize error manager for testing
	manager := errorCodes.GetManager()
	wd, _ := os.Getwd()
	// Navigate up from handler/store to project root
	errorFilePath := filepath.Join(wd, "..", "..", "..", "internal", "errors", "errors.yaml")
	if err := manager.LoadFromFile(errorFilePath); err != nil {
		panic("Failed to load errors.yaml for testing: " + err.Error())
	}

	// register custom validators for testing
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("taiwanlandline", utils.ValidateTaiwanLandline)
	}
}

func TestNewCreateStoreHandler(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockCreateStoreService is a mock implementation of CreateStoreServiceInterface
type MockCreateStoreService struct {
	mock.Mock
}

func (m *MockCreateStoreService) CreateStore(ctx context.Context, req store.CreateStoreRequest, staffContext common.StaffContext) (*store.CreateStoreResponse, error) {
	args := m.Called(ctx, req, staffContext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.CreateStoreResponse), args.Error(1)
}

func TestCreateStoreHandler_CreateStore_Success(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	createReq := store.CreateStoreRequest{
		Name:    "大安旗艦店",
		Address: "台北市大安區復興南路一段100號",
		Phone:   "02-12345678",
	}
	reqBody, _ := json.Marshal(createReq)

	c.Request, _ = http.NewRequest("POST", "/api/stores", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service response
	expectedResponse := &store.CreateStoreResponse{
		ID:       "8000000001",
		Name:     "大安旗艦店",
		Address:  "台北市大安區復興南路一段100號",
		Phone:    "02-12345678",
		IsActive: true,
	}
	mockService.On("CreateStore", mock.Anything, mock.AnythingOfType("store.CreateStoreRequest"), staffContext).Return(expectedResponse, nil)

	handler.CreateStore(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateStoreHandler_CreateStore_InvalidJSON(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with invalid JSON
	reqBody := `{"name":}`
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.CreateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateStoreHandler_CreateStore_ValidationError(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without required name field
	reqBody := `{"address":"台北市大安區復興南路一段100號","phone":"02-12345678"}`
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.CreateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateStoreHandler_CreateStore_MissingStaffContext(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request without staff context
	reqBody := `{"name":"測試店"}`
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	handler.CreateStore(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateStoreHandler_CreateStore_ServiceError(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request
	reqBody := `{"name":"重複店名"}`
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	// Mock service error
	serviceError := errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "store name already exists", nil)
	mockService.On("CreateStore", mock.Anything, mock.AnythingOfType("store.CreateStoreRequest"), staffContext).Return(nil, serviceError)

	handler.CreateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateStoreHandler_CreateStore_InvalidTaiwanLandline(t *testing.T) {
	mockService := &MockCreateStoreService{}
	handler := NewCreateStoreHandler(mockService)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up request with invalid Taiwan landline
	reqBody := `{"name":"測試店","address":"台北市大安區復興南路一段100號","phone":"invalid-phone"}`
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Set staff context
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}
	c.Set("user", staffContext)

	handler.CreateStore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}
