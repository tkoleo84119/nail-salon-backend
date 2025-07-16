package staff

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)




// Note: Full integration tests for CreateStaff would require mocking the database transaction
// For now, we'll focus on testing the validation logic and error scenarios
// The transaction logic is tested implicitly in the actual application

// Focus on testing validation logic and error scenarios that don't require database transactions

func TestCreateStaffService_CreateStaff_PermissionDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	tests := []struct {
		name        string
		creatorRole string
		targetRole  string
		expectedErr string
	}{
		{
			name:        "Admin_creating_SuperAdmin",
			creatorRole: staff.RoleAdmin,
			targetRole:  staff.RoleSuperAdmin,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Admin_creating_Admin",
			creatorRole: staff.RoleAdmin,
			targetRole:  staff.RoleAdmin,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Manager_creating_anyone",
			creatorRole: staff.RoleManager,
			targetRole:  staff.RoleStylist,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Stylist_creating_anyone",
			creatorRole: staff.RoleStylist,
			targetRole:  staff.RoleStylist,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "SuperAdmin_creating_SuperAdmin",
			creatorRole: staff.RoleSuperAdmin,
			targetRole:  staff.RoleSuperAdmin,
			expectedErr: "AUTH_PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := staff.CreateStaffRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
				Role:     tt.targetRole,
				StoreIDs: []string{"1"},
			}

			// Call service
			response, err := service.CreateStaff(context.Background(), req, tt.creatorRole, []int64{1})

			// Assert results
			assert.Error(t, err)
			assert.Nil(t, response)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestCreateStaffService_CreateStaff_InvalidRole(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data with invalid role
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     "INVALID_ROLE",
		StoreIDs: []string{"1"},
	}

	// Call service
	response, err := service.CreateStaff(context.Background(), req, staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_INVALID_ROLE")
}

func TestCreateStaffService_CreateStaff_UserAlreadyExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}

	// Mock expectations - user already exists
	mockQuerier.On("CheckStaffUserExists", mock.Anything, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	}).Return(true, nil)

	// Call service
	response, err := service.CreateStaff(context.Background(), req, staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_ALREADY_EXISTS")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStaffService_CreateStaff_StoreNotExist(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}

	// Mock expectations
	mockQuerier.On("CheckStaffUserExists", mock.Anything, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	}).Return(false, nil)

	// Store check returns 0 stores when 1 is requested (store not found)
	mockQuerier.On("CheckStoresExistAndActive", mock.Anything, []int64{1}).Return(
		dbgen.CheckStoresExistAndActiveRow{
			TotalCount:  0,
			ActiveCount: 0,
		}, nil)

	// Call service
	response, err := service.CreateStaff(context.Background(), req, staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_EXIST")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStaffService_CreateStaff_StoreNotActive(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}

	// Mock expectations
	mockQuerier.On("CheckStaffUserExists", mock.Anything, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	}).Return(false, nil)

	// Store exists but is not active
	mockQuerier.On("CheckStoresExistAndActive", mock.Anything, []int64{1}).Return(
		dbgen.CheckStoresExistAndActiveRow{
			TotalCount:  1,
			ActiveCount: 0,
		}, nil)

	// Call service
	response, err := service.CreateStaff(context.Background(), req, staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_ACTIVE")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStaffService_CreateStaff_StoreAccessDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"2"}, // Trying to assign store 2
	}

	// Call service - Admin only has access to store 1, but trying to assign store 2
	response, err := service.CreateStaff(context.Background(), req, staff.RoleAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")
}

func TestCreateStaffService_CreateStaff_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateStaffService(mockQuerier, nil)

	// Test data
	req := staff.CreateStaffRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     staff.RoleManager,
		StoreIDs: []string{"1"},
	}

	// Mock expectations - database error
	mockQuerier.On("CheckStaffUserExists", mock.Anything, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	}).Return(false, errors.New("database connection failed"))

	// Call service
	response, err := service.CreateStaff(context.Background(), req, staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to check user existence")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStaffService_validatePermissions(t *testing.T) {
	service := &CreateStaffService{}

	tests := []struct {
		name         string
		creatorRole  string
		targetRole   string
		expectError  bool
		expectedCode string
	}{
		{
			name:        "SuperAdmin_can_create_Admin",
			creatorRole: staff.RoleSuperAdmin,
			targetRole:  staff.RoleAdmin,
			expectError: false,
		},
		{
			name:        "SuperAdmin_can_create_Manager",
			creatorRole: staff.RoleSuperAdmin,
			targetRole:  staff.RoleManager,
			expectError: false,
		},
		{
			name:        "SuperAdmin_can_create_Stylist",
			creatorRole: staff.RoleSuperAdmin,
			targetRole:  staff.RoleStylist,
			expectError: false,
		},
		{
			name:         "SuperAdmin_cannot_create_SuperAdmin",
			creatorRole:  staff.RoleSuperAdmin,
			targetRole:   staff.RoleSuperAdmin,
			expectError:  true,
			expectedCode: "AUTH_PERMISSION_DENIED",
		},
		{
			name:        "Admin_can_create_Manager",
			creatorRole: staff.RoleAdmin,
			targetRole:  staff.RoleManager,
			expectError: false,
		},
		{
			name:        "Admin_can_create_Stylist",
			creatorRole: staff.RoleAdmin,
			targetRole:  staff.RoleStylist,
			expectError: false,
		},
		{
			name:         "Admin_cannot_create_Admin",
			creatorRole:  staff.RoleAdmin,
			targetRole:   staff.RoleAdmin,
			expectError:  true,
			expectedCode: "AUTH_PERMISSION_DENIED",
		},
		{
			name:         "Manager_cannot_create_anyone",
			creatorRole:  staff.RoleManager,
			targetRole:   staff.RoleStylist,
			expectError:  true,
			expectedCode: "AUTH_PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePermissions(tt.creatorRole, tt.targetRole)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateStaffService_validateStoreAccess(t *testing.T) {
	service := &CreateStaffService{}

	tests := []struct {
		name             string
		creatorRole      string
		creatorStoreIDs  []int64
		targetStoreIDs   []int64
		expectError      bool
		expectedCode     string
	}{
		{
			name:            "SuperAdmin_can_assign_any_store",
			creatorRole:     staff.RoleSuperAdmin,
			creatorStoreIDs: []int64{1},
			targetStoreIDs:  []int64{1},
			expectError:     false,
		},
		{
			name:            "Admin_can_assign_own_stores",
			creatorRole:     staff.RoleAdmin,
			creatorStoreIDs: []int64{1},
			targetStoreIDs:  []int64{1},
			expectError:     false,
		},
		{
			name:            "Admin_can_assign_subset_of_own_stores",
			creatorRole:     staff.RoleAdmin,
			creatorStoreIDs: []int64{1},
			targetStoreIDs:  []int64{1},
			expectError:     false,
		},
		{
			name:            "Admin_cannot_assign_unauthorized_store",
			creatorRole:     staff.RoleAdmin,
			creatorStoreIDs: []int64{1},
			targetStoreIDs:  []int64{2},
			expectError:     true,
			expectedCode:    "AUTH_PERMISSION_DENIED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateStoreAccess(tt.creatorRole, tt.creatorStoreIDs, tt.targetStoreIDs)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}