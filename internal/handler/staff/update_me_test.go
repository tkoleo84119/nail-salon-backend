package staff

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
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
)

// MockUpdateStaffMeService implements the UpdateStaffMeServiceInterface for testing
type MockUpdateStaffMeService struct {
	mock.Mock
}

// Ensure MockUpdateStaffMeService implements the interface
var _ staffService.UpdateStaffMeServiceInterface = (*MockUpdateStaffMeService)(nil)

func (m *MockUpdateStaffMeService) UpdateStaffMe(ctx context.Context, req staff.UpdateStaffMeRequest, staffUserID int64) (*staff.UpdateStaffMeResponse, error) {
	args := m.Called(ctx, req, staffUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.UpdateStaffMeResponse), args.Error(1)
}

func setupTestGinForUpdateMe() {
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	_ = errorManager.LoadFromFile("../../errors/errors.yaml")
}

func TestUpdateStaffMeHandler_UpdateStaffMe_Success(t *testing.T) {
	setupTestGinForUpdateMe()

	mockService := new(MockUpdateStaffMeService)
	handler := NewUpdateStaffMeHandler(mockService)

	reqBody := staff.UpdateStaffMeRequest{
		Email: stringPtr("new-email@example.com"),
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedResponse := &staff.UpdateStaffMeResponse{
		ID:       "12345",
		Username: "staff_amy",
		Email:    "new-email@example.com",
		Role:     staff.RoleAdmin,
	}

	mockService.On("UpdateStaffMe", mock.Anything, reqBody, int64(12345)).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/staff/me", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	handler.UpdateStaffMe(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestUpdateStaffMeHandler_UpdateStaffMe_MissingStaffUserID(t *testing.T) {
	setupTestGinForUpdateMe()

	mockService := new(MockUpdateStaffMeService)
	handler := NewUpdateStaffMeHandler(mockService)

	reqBody := staff.UpdateStaffMeRequest{
		Email: stringPtr("new-email@example.com"),
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/staff/me", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	// Missing staff_user_id in context

	handler.UpdateStaffMe(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateStaffMeHandler_UpdateStaffMe_ServiceError(t *testing.T) {
	setupTestGinForUpdateMe()

	mockService := new(MockUpdateStaffMeService)
	handler := NewUpdateStaffMeHandler(mockService)

	reqBody := staff.UpdateStaffMeRequest{
		Email: stringPtr("existing-email@example.com"),
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("UpdateStaffMe", mock.Anything, reqBody, int64(12345)).Return(nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserEmailExists))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/api/staff/me", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set staff context
	staffContext := common.StaffContext{
		UserID:   "12345",
		Username: "admin",
		Role:     staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "1", Name: "Test Store"},
		},
	}
	c.Set("user", staffContext)

	handler.UpdateStaffMe(c)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockService.AssertExpectations(t)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
