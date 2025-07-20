package service

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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
)

// MockCreateServiceService implements the CreateServiceInterface for testing
type MockCreateServiceService struct {
	mock.Mock
}

// Ensure MockCreateServiceService implements the interface
var _ serviceService.CreateServiceInterface = (*MockCreateServiceService)(nil)

func (m *MockCreateServiceService) CreateService(ctx context.Context, req service.CreateServiceRequest, creatorRole string) (*service.CreateServiceResponse, error) {
	args := m.Called(ctx, req, creatorRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.CreateServiceResponse), args.Error(1)
}

func setupTestGinForCreate() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestCreateServiceHandler_CreateService_Success(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateServiceService)
	handler := NewCreateServiceHandler(mockService)

	// Test request
	req := service.CreateServiceRequest{
		Name:            "凝膠手部單色",
		Price:           1200,
		DurationMinutes: 60,
		IsAddon:         false,
		IsVisible:       true,
		Note:            "含基礎修型保養",
	}

	// Expected response from service
	expectedResponse := &service.CreateServiceResponse{
		ID:              "9000000001",
		Name:            "凝膠手部單色",
		Price:           1200,
		DurationMinutes: 60,
		IsAddon:         false,
		IsVisible:       true,
		IsActive:        true,
		Note:            "含基礎修型保養",
	}

	// Setup mock expectations
	mockService.On("CreateService", mock.Anything, req, staff.RoleAdmin).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context (simulate authenticated admin user)
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Main Store"},
		},
	}
	c.Set("user", staffContext)

	// Call the handler
	handler.CreateService(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Verify the response data
	responseData, _ := json.Marshal(response.Data)
	var actualResponse service.CreateServiceResponse
	json.Unmarshal(responseData, &actualResponse)

	assert.Equal(t, expectedResponse.ID, actualResponse.ID)
	assert.Equal(t, expectedResponse.Name, actualResponse.Name)
	assert.Equal(t, expectedResponse.Price, actualResponse.Price)
	assert.Equal(t, expectedResponse.DurationMinutes, actualResponse.DurationMinutes)
	assert.Equal(t, expectedResponse.IsAddon, actualResponse.IsAddon)
	assert.Equal(t, expectedResponse.IsVisible, actualResponse.IsVisible)
	assert.Equal(t, expectedResponse.IsActive, actualResponse.IsActive)
	assert.Equal(t, expectedResponse.Note, actualResponse.Note)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestCreateServiceHandler_CreateService_ValidationErrors(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateServiceService)
	handler := NewCreateServiceHandler(mockService)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
	}{
		{
			name: "Empty name",
			request: service.CreateServiceRequest{
				Name:            "",
				Price:           1200,
				DurationMinutes: 60,
				IsAddon:         false,
				IsVisible:       true,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Zero price",
			request: service.CreateServiceRequest{
				Name:            "Test Service",
				Price:           0,
				DurationMinutes: 60,
				IsAddon:         false,
				IsVisible:       true,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid duration (too high)",
			request: service.CreateServiceRequest{
				Name:            "Test Service",
				Price:           1200,
				DurationMinutes: 1500,
				IsAddon:         false,
				IsVisible:       true,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Note too long",
			request: service.CreateServiceRequest{
				Name:            "Test Service",
				Price:           1200,
				DurationMinutes: 60,
				IsAddon:         false,
				IsVisible:       true,
				Note:            string(make([]byte, 300)), // 300 characters, exceeds 255 limit
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			jsonData, _ := json.Marshal(tt.request)
			httpReq, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(jsonData))
			httpReq.Header.Set("Content-Type", "application/json")

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httpReq

			// Set staff context
			staffContext := common.StaffContext{
				UserID:   "1",
				Username: "admin",
				Role:     staff.RoleAdmin,
				StoreList: []common.Store{
					{ID: "1", Name: "Main Store"},
				},
			}
			c.Set("user", staffContext)

			// Call the handler
			handler.CreateService(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateServiceHandler_CreateService_NoStaffContext(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateServiceService)
	handler := NewCreateServiceHandler(mockService)

	// Test request
	req := service.CreateServiceRequest{
		Name:            "Test Service",
		Price:           1200,
		DurationMinutes: 60,
		IsAddon:         false,
		IsVisible:       true,
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context without staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Call the handler
	handler.CreateService(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify mock was not called
	mockService.AssertNotCalled(t, "CreateService")
}

func TestCreateServiceHandler_CreateService_ServiceError(t *testing.T) {
	setupTestGinForCreate()

	// Create mock service
	mockService := new(MockCreateServiceService)
	handler := NewCreateServiceHandler(mockService)

	// Test request
	req := service.CreateServiceRequest{
		Name:            "Test Service",
		Price:           1200,
		DurationMinutes: 60,
		IsAddon:         false,
		IsVisible:       true,
	}

	// Setup mock to return error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.UserAlreadyExists)
	mockService.On("CreateService", mock.Anything, req, staff.RoleAdmin).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/api/services", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "1",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Main Store"},
		},
	}
	c.Set("user", staffContext)

	// Call the handler
	handler.CreateService(c)

	// Assertions
	assert.Equal(t, http.StatusConflict, w.Code)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}