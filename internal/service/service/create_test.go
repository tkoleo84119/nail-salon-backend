package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestCreateServiceService_CreateService_PermissionDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	serviceService := NewCreateServiceService(mockQuerier)

	// Test request
	req := service.CreateServiceRequest{
		Name:            "Test Service",
		Price:           1200,
		DurationMinutes: 60,
		IsAddon:         false,
		IsVisible:       true,
	}

	tests := []struct {
		name        string
		creatorRole string
		expectedErr string
	}{
		{
			name:        "Manager cannot create service",
			creatorRole: staff.RoleManager,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Stylist cannot create service",
			creatorRole: staff.RoleStylist,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Invalid role cannot create service",
			creatorRole: "INVALID_ROLE",
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the service
			response, err := serviceService.CreateService(context.Background(), req, tt.creatorRole)

			// Assertions
			assert.Error(t, err)
			assert.Nil(t, response)

			if serviceErr, ok := err.(*errorCodes.ServiceError); ok {
				assert.Equal(t, tt.expectedErr, serviceErr.Code)
			}
		})
	}

	// Verify no querier calls were made
	mockQuerier.AssertNotCalled(t, "GetServiceByName")
	mockQuerier.AssertNotCalled(t, "CreateService")
}

func TestCreateServiceService_ValidatePermissions(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	serviceService := NewCreateServiceService(mockQuerier)

	tests := []struct {
		name        string
		creatorRole string
		expectError bool
	}{
		{
			name:        "SuperAdmin can create service",
			creatorRole: staff.RoleSuperAdmin,
			expectError: false,
		},
		{
			name:        "Admin can create service",
			creatorRole: staff.RoleAdmin,
			expectError: false,
		},
		{
			name:        "Manager cannot create service",
			creatorRole: staff.RoleManager,
			expectError: true,
		},
		{
			name:        "Stylist cannot create service",
			creatorRole: staff.RoleStylist,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := serviceService.validatePermissions(tt.creatorRole)

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