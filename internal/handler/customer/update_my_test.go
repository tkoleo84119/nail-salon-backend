package customer

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

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

// MockUpdateMyCustomerService for testing
type MockUpdateMyCustomerService struct {
	mock.Mock
}

func (m *MockUpdateMyCustomerService) UpdateMyCustomer(ctx context.Context, customerID int64, req customer.UpdateMyCustomerRequest) (*customer.UpdateMyCustomerResponse, error) {
	args := m.Called(ctx, customerID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*customer.UpdateMyCustomerResponse), args.Error(1)
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	customerID := int64(1000000001)
	name := "王小美"
	phone := "0912345678"
	birthday := "1992-02-29"
	city := "台北市"
	favoriteShapes := []string{"方形"}
	favoriteColors := []string{"粉色"}
	favoriteStyles := []string{"法式"}
	isIntrovert := true
	customerNote := "容易指緣乾裂"

	requestBody := map[string]interface{}{
		"name":           name,
		"phone":          phone,
		"birthday":       birthday,
		"city":           city,
		"favoriteShapes": favoriteShapes,
		"favoriteColors": favoriteColors,
		"favoriteStyles": favoriteStyles,
		"isIntrovert":    isIntrovert,
		"customerNote":   customerNote,
	}

	expectedRequest := customer.UpdateMyCustomerRequest{
		Name:           &name,
		Phone:          &phone,
		Birthday:       &birthday,
		City:           &city,
		FavoriteShapes: &favoriteShapes,
		FavoriteColors: &favoriteColors,
		FavoriteStyles: &favoriteStyles,
		IsIntrovert:    &isIntrovert,
		CustomerNote:   &customerNote,
	}

	expectedResponse := &customer.UpdateMyCustomerResponse{
		ID:             "1000000001",
		Name:           name,
		Phone:          phone,
		Birthday:       &birthday,
		City:           &city,
		FavoriteShapes: &favoriteShapes,
		FavoriteColors: &favoriteColors,
		FavoriteStyles: &favoriteStyles,
		IsIntrovert:    &isIntrovert,
		CustomerNote:   &customerNote,
	}

	mockService.On("UpdateMyCustomer", mock.Anything, customerID, expectedRequest).Return(expectedResponse, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		customerContext := common.CustomerContext{
			CustomerID: customerID,
		}
		c.Set(middleware.CustomerContextKey, customerContext)
		c.Next()
	})
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	responseData := response.Data.(map[string]interface{})
	assert.Equal(t, expectedResponse.ID, responseData["id"])
	assert.Equal(t, expectedResponse.Name, responseData["name"])
	assert.Equal(t, expectedResponse.Phone, responseData["phone"])

	mockService.AssertExpectations(t)
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_MissingCustomerContext(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	requestBody := map[string]interface{}{
		"name": "王小美",
	}

	router := gin.New()
	// No customer context set
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "未找到使用者認證資訊")

	mockService.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_InvalidJSON(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	customerID := int64(1000000001)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		customerContext := common.CustomerContext{
			CustomerID: customerID,
		}
		c.Set(middleware.CustomerContextKey, customerContext)
		c.Next()
	})
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	invalidJSON := `{"name": "王小美", "phone": "invalid phone format}`
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_NoFieldsToUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	customerID := int64(1000000001)

	requestBody := map[string]interface{}{} // Empty request

	router := gin.New()
	router.Use(func(c *gin.Context) {
		customerContext := common.CustomerContext{
			CustomerID: customerID,
		}
		c.Set(middleware.CustomerContextKey, customerContext)
		c.Next()
	})
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "至少需要提供一個欄位進行更新")

	mockService.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_ServiceError(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	customerID := int64(1000000001)
	name := "王小美"

	requestBody := map[string]interface{}{
		"name": name,
	}

	expectedRequest := customer.UpdateMyCustomerRequest{
		Name: &name,
	}

	mockService.On("UpdateMyCustomer", mock.Anything, customerID, expectedRequest).Return(nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound))

	router := gin.New()
	router.Use(func(c *gin.Context) {
		customerContext := common.CustomerContext{
			CustomerID: customerID,
		}
		c.Set(middleware.CustomerContextKey, customerContext)
		c.Next()
	})
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_ValidationErrors(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	testCases := []struct {
		name        string
		requestBody map[string]interface{}
		expectError string
	}{
		{
			name: "Invalid phone format",
			requestBody: map[string]interface{}{
				"phone": "invalid-phone",
			},
			expectError: "phone",
		},
		{
			name: "Name too long",
			requestBody: map[string]interface{}{
				"name": string(make([]byte, 101)), // 101 characters, max is 100
			},
			expectError: "name",
		},
		{
			name: "Customer note too long",
			requestBody: map[string]interface{}{
				"customerNote": string(make([]byte, 1001)), // 1001 characters, max is 1000
			},
			expectError: "customerNote",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockUpdateMyCustomerService{}
			handler := NewUpdateMyCustomerHandler(mockService)

			customerID := int64(1000000001)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				customerContext := common.CustomerContext{
					CustomerID: customerID,
				}
				c.Set(middleware.CustomerContextKey, customerContext)
				c.Next()
			})
			router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			if tc.expectError != "" {
				assert.Contains(t, w.Body.String(), tc.expectError)
			}

			mockService.AssertNotCalled(t, "UpdateMyCustomer")
		})
	}
}

func TestUpdateMyCustomerHandler_UpdateMyCustomer_OnlyNameUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForHandler(t)
	defer env.Cleanup()

	mockService := &MockUpdateMyCustomerService{}
	handler := NewUpdateMyCustomerHandler(mockService)

	customerID := int64(1000000001)
	name := "新名字"

	requestBody := map[string]interface{}{
		"name": name,
	}

	expectedRequest := customer.UpdateMyCustomerRequest{
		Name: &name,
	}

	expectedResponse := &customer.UpdateMyCustomerResponse{
		ID:    "1000000001",
		Name:  name,
		Phone: "0912345678",
	}

	mockService.On("UpdateMyCustomer", mock.Anything, customerID, expectedRequest).Return(expectedResponse, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		customerContext := common.CustomerContext{
			CustomerID: customerID,
		}
		c.Set(middleware.CustomerContextKey, customerContext)
		c.Next()
	})
	router.PATCH("/api/customers/me", handler.UpdateMyCustomer)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("PATCH", "/api/customers/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response common.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	responseData := response.Data.(map[string]interface{})
	assert.Equal(t, expectedResponse.Name, responseData["name"])

	mockService.AssertExpectations(t)
}