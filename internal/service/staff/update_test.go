package staff

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockUpdateStaffRepository mocks the sqlx repository for testing
type MockUpdateStaffRepository struct {
	mock.Mock
}

// Ensure MockUpdateStaffRepository implements the interface
var _ sqlxRepo.StaffUserRepositoryInterface = (*MockUpdateStaffRepository)(nil)

func (m *MockUpdateStaffRepository) UpdateStaffUser(ctx context.Context, id int64, req staff.UpdateStaffRequest) (*staff.UpdateStaffResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.UpdateStaffResponse), args.Error(1)
}

// MockUpdateStaffQuerier extends MockQuerier for update staff specific queries
type MockUpdateStaffQuerier struct {
	MockQuerier
}

func (m *MockUpdateStaffQuerier) GetStaffUserByID(ctx context.Context, id int64) (dbgen.StaffUser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func setupTestEnvironmentForUpdate(t *testing.T) func() {
	// Initialize snowflake for testing
	err := utils.InitSnowflake(1)
	require.NoError(t, err)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	err = errorManager.LoadFromFile("../../errors/errors.yaml")
	require.NoError(t, err)

	return func() {
		// cleanup if needed
	}
}

func TestUpdateStaffService_UpdateStaff_Success(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	isActive := false
	req := staff.UpdateStaffRequest{
		Role:     &role,
		IsActive: &isActive,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	expectedResponse := &staff.UpdateStaffResponse{
		ID:       "123456789",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: false,
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)
	mockRepo.On("UpdateStaffUser", mock.Anything, int64(123456789), req).Return(expectedResponse, nil)

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Username, response.Username)
	assert.Equal(t, expectedResponse.Email, response.Email)
	assert.Equal(t, expectedResponse.Role, response.Role)
	assert.Equal(t, expectedResponse.IsActive, response.IsActive)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_InvalidID(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Call service with invalid ID
	response, err := service.UpdateStaff(context.Background(), "invalid_id", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "VAL_INPUT_VALIDATION_FAILED")
}

func TestUpdateStaffService_UpdateStaff_EmptyRequest(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Empty request
	req := staff.UpdateStaffRequest{}

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "VAL_ALL_FIELDS_EMPTY")
}

func TestUpdateStaffService_UpdateStaff_InvalidRole(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data with invalid role
	invalidRole := "INVALID_ROLE"
	req := staff.UpdateStaffRequest{
		Role: &invalidRole,
	}

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_INVALID_ROLE")
}

func TestUpdateStaffService_UpdateStaff_UserNotFound(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Mock expectations - user not found
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, sql.ErrNoRows)

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_FOUND")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_CannotUpdateSelf(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       1, // Same as updater ID
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(1)).Return(*targetStaff, nil)

	// Call service - updater ID is same as target ID
	response, err := service.UpdateStaff(context.Background(), "1", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_UPDATE_SELF")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_CannotUpdateSuperAdmin(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "superadmin",
		Email:    "superadmin@example.com",
		Role:     staff.RoleSuperAdmin, // Target is SUPER_ADMIN
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_ManagerCannotUpdate(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleStylist
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)

	// Call service - updater is MANAGER (not allowed to update)
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleManager)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_AdminCannotSetToSuperAdmin(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data - trying to set to SUPER_ADMIN
	role := staff.RoleSuperAdmin
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)

	// Call service - updater is ADMIN trying to set to SUPER_ADMIN
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_AdminCannotSetToAdmin(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data - trying to set to ADMIN
	role := staff.RoleAdmin
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)

	// Call service - updater is ADMIN trying to set to ADMIN
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_AdminCanUpdateManagerAndStylist(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	tests := []struct {
		name       string
		targetRole string
	}{
		{
			name:       "Admin_can_update_to_Manager",
			targetRole: staff.RoleManager,
		},
		{
			name:       "Admin_can_update_to_Stylist",
			targetRole: staff.RoleStylist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockRepo.ExpectedCalls = nil

			// Test data
			req := staff.UpdateStaffRequest{
				Role: &tt.targetRole,
			}

			targetStaff := &dbgen.StaffUser{
				ID:       123456789,
				Username: "testuser",
				Email:    "test@example.com",
				Role:     staff.RoleStylist,
				IsActive: pgtype.Bool{Bool: true, Valid: true},
			}

			expectedResponse := &staff.UpdateStaffResponse{
				ID:       "123456789",
				Username: "testuser",
				Email:    "test@example.com",
				Role:     tt.targetRole,
				IsActive: true,
			}

			// Mock expectations
			mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)
			mockRepo.On("UpdateStaffUser", mock.Anything, int64(123456789), req).Return(expectedResponse, nil)

			// Call service - updater is ADMIN
			response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleAdmin)

			// Assert results
			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, expectedResponse.Role, response.Role)

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateStaffService_UpdateStaff_DatabaseError(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	// Mock expectations - database error
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, errors.New("database connection failed"))

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get target staff")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffService_UpdateStaff_UpdateDatabaseError(t *testing.T) {
	cleanup := setupTestEnvironmentForUpdate(t)
	defer cleanup()

	// Create mock repository and querier
	mockRepo := new(MockUpdateStaffRepository)
	mockQuerier := new(MockUpdateStaffQuerier)
	service := &UpdateStaffService{queries: mockQuerier, repo: mockRepo}

	// Test data
	role := staff.RoleManager
	req := staff.UpdateStaffRequest{
		Role: &role,
	}

	targetStaff := &dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(*targetStaff, nil)
	mockRepo.On("UpdateStaffUser", mock.Anything, int64(123456789), req).Return(nil, errors.New("update failed"))

	// Call service
	response, err := service.UpdateStaff(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to update staff")

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQuerier.AssertExpectations(t)
}