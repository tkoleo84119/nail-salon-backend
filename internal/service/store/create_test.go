package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

// Note: Full integration tests for CreateStore would require mocking the database transaction
// For now, we'll focus on testing the validation logic and error scenarios
// The transaction logic is tested implicitly in the actual application

func TestCreateStoreService_CreateStore_InvalidUserID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStoreService(mockQuerier, nil)

	ctx := context.Background()
	req := store.CreateStoreRequest{
		Name: "測試店",
	}
	staffContext := common.StaffContext{
		UserID: "invalid",
		Role:   staff.RoleAdmin,
	}

	response, err := service.CreateStore(ctx, req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreService_CreateStore_InsufficientPermission(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		shouldError bool
	}{
		{
			name:        "Manager_cannot_create_store",
			role:        staff.RoleManager,
			shouldError: true,
		},
		{
			name:        "Stylist_cannot_create_store",
			role:        staff.RoleStylist,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuerier := mocks.NewMockQuerier()
			service := NewCreateStoreService(mockQuerier, nil)
			
			req := store.CreateStoreRequest{
				Name: "測試店",
			}
			staffContext := common.StaffContext{
				UserID: "33333",
				Role:   tt.role,
			}

			response, err := service.CreateStore(context.Background(), req, staffContext)

			assert.Nil(t, response)
			assert.Error(t, err)
			serviceErr, ok := err.(*errorCodes.ServiceError)
			assert.True(t, ok)
			assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

			mockQuerier.AssertExpectations(t)
		})
	}
}

func TestCreateStoreService_CreateStore_NameAlreadyExists(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStoreService(mockQuerier, nil)

	ctx := context.Background()
	req := store.CreateStoreRequest{
		Name: "重複店名",
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock name already exists
	mockQuerier.On("CheckStoreNameExists", ctx, req.Name).Return(true, nil)

	response, err := service.CreateStore(ctx, req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreService_CreateStore_DatabaseError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStoreService(mockQuerier, nil)

	ctx := context.Background()
	req := store.CreateStoreRequest{
		Name: "測試店",
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock database error when checking name
	mockQuerier.On("CheckStoreNameExists", ctx, req.Name).Return(false, errors.New("database error"))

	response, err := service.CreateStore(ctx, req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}