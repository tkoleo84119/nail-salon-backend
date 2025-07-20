package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestUpdateServiceService_UpdateService_PermissionDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	tests := []struct {
		name        string
		updaterRole string
		expectedErr string
	}{
		{
			name:        "Manager cannot update service",
			updaterRole: staff.RoleManager,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Stylist cannot update service",
			updaterRole: staff.RoleStylist,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Invalid role cannot update service",
			updaterRole: "INVALID_ROLE",
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the service
			response, err := serviceService.UpdateService(context.Background(), "123", req, tt.updaterRole)

			// Assertions
			assert.Error(t, err)
			assert.Nil(t, response)

			if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
				assert.Equal(t, tt.expectedErr, serviceErr.Code)
			}
		})
	}

	// Verify no querier or repository calls were made
	mockQuerier.AssertNotCalled(t, "GetServiceByID")
	mockQuerier.AssertNotCalled(t, "CheckServiceNameExistsExcluding")
	mockRepo.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceService_UpdateService_InvalidServiceID(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Call the service with invalid ID
	response, err := serviceService.UpdateService(context.Background(), "invalid", req, staff.RoleAdmin)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "VAL_INPUT_VALIDATION_FAILED", serviceErr.Code)
	}

	// Verify no calls were made
	mockQuerier.AssertNotCalled(t, "GetServiceByID")
	mockQuerier.AssertNotCalled(t, "CheckServiceNameExistsExcluding")
	mockRepo.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceService_UpdateService_NoFieldsToUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	// Empty request
	req := service.UpdateServiceRequest{}

	// Call the service
	response, err := serviceService.UpdateService(context.Background(), "123", req, staff.RoleAdmin)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "VAL_ALL_FIELDS_EMPTY", serviceErr.Code)
	}

	// Verify no calls were made
	mockQuerier.AssertNotCalled(t, "GetServiceByID")
	mockQuerier.AssertNotCalled(t, "CheckServiceNameExistsExcluding")
	mockRepo.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceService_UpdateService_ServiceNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	// Test request
	nameUpdate := "Test Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Mock service not found
	mockQuerier.On("GetServiceByID", mock.Anything, int64(123)).Return(dbgen.Service{}, pgx.ErrNoRows)

	// Call the service
	response, err := serviceService.UpdateService(context.Background(), "123", req, staff.RoleAdmin)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "SERVICE_NOT_FOUND", serviceErr.Code)
	}

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceService_UpdateService_ServiceNameExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	// Test request
	nameUpdate := "Existing Service"
	req := service.UpdateServiceRequest{
		Name: &nameUpdate,
	}

	// Mock existing service
	existingService := dbgen.Service{
		ID:   123,
		Name: "Current Service",
	}
	mockQuerier.On("GetServiceByID", mock.Anything, int64(123)).Return(existingService, nil)

	// Mock name already exists
	mockQuerier.On("CheckServiceNameExistsExcluding", mock.Anything, dbgen.CheckServiceNameExistsExcludingParams{
		Name: "Existing Service",
		ID:   123,
	}).Return(true, nil)

	// Call the service
	response, err := serviceService.UpdateService(context.Background(), "123", req, staff.RoleAdmin)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)

	if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
		assert.Equal(t, "SERVICE_NAME_ALREADY_EXISTS", serviceErr.Code)
	}

	// Verify mock expectations
	mockQuerier.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "UpdateService")
}

func TestUpdateServiceService_ValidatePermissions(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier and repository
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockServiceRepository()
	serviceService := NewUpdateServiceService(mockQuerier, mockRepo)

	tests := []struct {
		name        string
		updaterRole string
		expectError bool
	}{
		{
			name:        "SuperAdmin can update service",
			updaterRole: staff.RoleSuperAdmin,
			expectError: false,
		},
		{
			name:        "Admin can update service",
			updaterRole: staff.RoleAdmin,
			expectError: false,
		},
		{
			name:        "Manager cannot update service",
			updaterRole: staff.RoleManager,
			expectError: true,
		},
		{
			name:        "Stylist cannot update service",
			updaterRole: staff.RoleStylist,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := serviceService.validatePermissions(tt.updaterRole)

			if tt.expectError {
				assert.Error(t, err)
				if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
					assert.Equal(t, "AUTH_PERMISSION_DENIED", serviceErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}