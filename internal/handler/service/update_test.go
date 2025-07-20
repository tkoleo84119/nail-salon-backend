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

// MockUpdateServiceService implements the UpdateServiceInterface for testing
type MockUpdateServiceService struct {
	mock.Mock
}

// Ensure MockUpdateServiceService implements the interface
var _ serviceService.UpdateServiceInterface = (*MockUpdateServiceService)(nil)

func (m *MockUpdateServiceService) UpdateService(ctx context.Context, serviceID string, req service.UpdateServiceRequest, updaterRole string) (*service.UpdateServiceResponse, error) {
	args := m.Called(ctx, serviceID, req, updaterRole)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UpdateServiceResponse), args.Error(1)
}

func setupTestGinForUpdate() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestUpdateServiceHandler_UpdateService_Success(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	// Test request
	nameUpdate := "凝膠足部單色"
	priceUpdate := int64(1400)
	req := service.UpdateServiceRequest{
		Name:  &nameUpdate,
		Price: &priceUpdate,
	}

	// Expected response from service
	expectedResponse := &service.UpdateServiceResponse{
		ID:              "9000000001",
		Name:            "凝膠足部單色",
		Price:           1400,
		DurationMinutes: 75,
		IsAddon:         false,
		IsVisible:       true,
		IsActive:        true,
		Note:            "足部基礎保養",
	}

	// Setup mock expectations
	mockService.On("UpdateService", mock.Anything, "9000000001", req, staff.RoleAdmin).Return(expectedResponse, nil)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PATCH", "/api/services/9000000001", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "serviceId", Value: "9000000001"}}

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
	handler.UpdateService(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)

	// Verify the response data
	responseData, _ := json.Marshal(response.Data)
	var actualResponse service.UpdateServiceResponse
	json.Unmarshal(responseData, &actualResponse)

	assert.Equal(t, expectedResponse.ID, actualResponse.ID)
	assert.Equal(t, expectedResponse.Name, actualResponse.Name)
	assert.Equal(t, expectedResponse.Price, actualResponse.Price)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestUpdateServiceHandler_UpdateService_ValidationErrors(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
	}{
		{
			name: "Name too short",
			request: service.UpdateServiceRequest{
				Name: func() *string { s := ""; return &s }(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative price",
			request: service.UpdateServiceRequest{
				Price: func() *int64 { p := int64(-1); return &p }(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Duration too high",
			request: service.UpdateServiceRequest{
				DurationMinutes: func() *int32 { d := int32(1500); return &d }(),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Note too long",
			request: service.UpdateServiceRequest{
				Note: func() *string { s := string(make([]byte, 300)); return &s }(),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			jsonData, _ := json.Marshal(tt.request)
			httpReq, _ := http.NewRequest("PATCH", "/api/services/9000000001", bytes.NewBuffer(jsonData))
			httpReq.Header.Set("Content-Type", "application/json")

			// Create gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httpReq
			c.Params = gin.Params{{Key: "serviceId", Value: "9000000001"}}

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
			handler.UpdateService(c)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateServiceHandler_UpdateService_NoServiceId(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Create HTTP request without serviceId parameter
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PATCH", "/api/services/", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context without serviceId parameter
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
	handler.UpdateService(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify mock was not called
	mockService.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceHandler_UpdateService_NoFieldsToUpdate(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	// Empty request
	req := service.UpdateServiceRequest{}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PATCH", "/api/services/9000000001", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "serviceId", Value: "9000000001"}}

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
	handler.UpdateService(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify mock was not called
	mockService.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceHandler_UpdateService_NoStaffContext(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PATCH", "/api/services/9000000001", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context without staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "serviceId", Value: "9000000001"}}

	// Call the handler
	handler.UpdateService(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Verify mock was not called
	mockService.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceHandler_UpdateService_ServiceError(t *testing.T) {
	setupTestGinForUpdate()

	// Create mock service
	mockService := new(MockUpdateServiceService)
	handler := NewUpdateServiceHandler(mockService)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Setup mock to return error
	serviceError := errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
	mockService.On("UpdateService", mock.Anything, "9000000001", req, staff.RoleAdmin).Return(nil, serviceError)

	// Create HTTP request
	jsonData, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PATCH", "/api/services/9000000001", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create gin context with staff context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Params = gin.Params{{Key: "serviceId", Value: "9000000001"}}

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
	handler.UpdateService(c)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}