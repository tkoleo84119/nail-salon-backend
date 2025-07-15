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
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

// MockCreateStoreAccessQuerier mocks the querier for create store access
type MockCreateStoreAccessQuerier struct {
	MockQuerier
}

func (m *MockCreateStoreAccessQuerier) GetStaffUserByID(ctx context.Context, id int64) (dbgen.StaffUser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockCreateStoreAccessQuerier) GetStoreByID(ctx context.Context, id int64) (dbgen.GetStoreByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.GetStoreByIDRow), args.Error(1)
}

func (m *MockCreateStoreAccessQuerier) CheckStoreAccessExists(ctx context.Context, arg dbgen.CheckStoreAccessExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockCreateStoreAccessQuerier) CreateStaffUserStoreAccess(ctx context.Context, arg dbgen.CreateStaffUserStoreAccessParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockCreateStoreAccessQuerier) GetStaffUserStoreAccess(ctx context.Context, staffUserID int64) ([]dbgen.GetStaffUserStoreAccessRow, error) {
	args := m.Called(ctx, staffUserID)
	return args.Get(0).([]dbgen.GetStaffUserStoreAccessRow), args.Error(1)
}

func setupTestEnvironmentForStoreAccess(t *testing.T) func() {
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

func TestCreateStoreAccessService_CreateStoreAccess_Success_NewlyCreated(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	store := dbgen.GetStoreByIDRow{
		ID:       2,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	existsResult := false

	storeAccessList := []dbgen.GetStaffUserStoreAccessRow{
		{StoreID: 1, StoreName: "Store 1"},
		{StoreID: 2, StoreName: "Test Store"},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)
	mockQuerier.On("GetStoreByID", mock.Anything, int64(2)).Return(store, nil)
	mockQuerier.On("CheckStoreAccessExists", mock.Anything, dbgen.CheckStoreAccessExistsParams{
		StaffUserID: 123456789,
		StoreID:     2,
	}).Return(existsResult, nil)
	mockQuerier.On("CreateStaffUserStoreAccess", mock.Anything, dbgen.CreateStaffUserStoreAccessParams{
		StoreID:     2,
		StaffUserID: 123456789,
	}).Return(nil)
	mockQuerier.On("GetStaffUserStoreAccess", mock.Anything, int64(123456789)).Return(storeAccessList, nil)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.NoError(t, err)
	assert.True(t, isNewlyCreated)
	assert.NotNil(t, response)
	assert.Equal(t, "123456789", response.StaffUserID)
	assert.Len(t, response.StoreList, 2)
	assert.Equal(t, int64(1), response.StoreList[0].ID)
	assert.Equal(t, "Store 1", response.StoreList[0].Name)
	assert.Equal(t, int64(2), response.StoreList[1].ID)
	assert.Equal(t, "Test Store", response.StoreList[1].Name)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_Success_AlreadyExists(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	store := dbgen.GetStoreByIDRow{
		ID:       2,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	existsResult := true

	storeAccessList := []dbgen.GetStaffUserStoreAccessRow{
		{StoreID: 1, StoreName: "Store 1"},
		{StoreID: 2, StoreName: "Test Store"},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)
	mockQuerier.On("GetStoreByID", mock.Anything, int64(2)).Return(store, nil)
	mockQuerier.On("CheckStoreAccessExists", mock.Anything, dbgen.CheckStoreAccessExistsParams{
		StaffUserID: 123456789,
		StoreID:     2,
	}).Return(existsResult, nil)
	mockQuerier.On("GetStaffUserStoreAccess", mock.Anything, int64(123456789)).Return(storeAccessList, nil)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.NoError(t, err)
	assert.False(t, isNewlyCreated)
	assert.NotNil(t, response)
	assert.Equal(t, "123456789", response.StaffUserID)
	assert.Len(t, response.StoreList, 2)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_InvalidID(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	// Call service with invalid ID
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "invalid_id", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "VAL_INPUT_VALIDATION_FAILED")
}

func TestCreateStoreAccessService_CreateStoreAccess_StaffNotFound(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	// Mock expectations - staff not found
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, sql.ErrNoRows)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STAFF_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_CannotUpdateSelf(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       1, // Same as creator ID
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(1)).Return(targetStaff, nil)

	// Call service - creator ID same as target ID
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "1", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_UPDATE_SELF")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_CannotUpdateSuperAdmin(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "superadmin",
		Email:    "superadmin@example.com",
		Role:     staff.RoleSuperAdmin, // Target is SUPER_ADMIN
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_StoreNotFound(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)
	mockQuerier.On("GetStoreByID", mock.Anything, int64(2)).Return(dbgen.GetStoreByIDRow{}, sql.ErrNoRows)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_StoreNotActive(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	store := dbgen.GetStoreByIDRow{
		ID:       2,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: false, Valid: true}, // Store is not active
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)
	mockQuerier.On("GetStoreByID", mock.Anything, int64(2)).Return(store, nil)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_ACTIVE")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_AdminNoPermission(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 3, // Admin doesn't have access to store 3
	}

	targetStaff := dbgen.StaffUser{
		ID:       123456789,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     staff.RoleManager,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	store := dbgen.GetStoreByIDRow{
		ID:       3,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(targetStaff, nil)
	mockQuerier.On("GetStoreByID", mock.Anything, int64(3)).Return(store, nil)

	// Call service - admin only has access to stores 1 and 2, but requesting store 3
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_DatabaseError(t *testing.T) {
	cleanup := setupTestEnvironmentForStoreAccess(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockCreateStoreAccessQuerier)
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: 2,
	}

	// Mock expectations - database error
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, errors.New("database connection failed"))

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1, 2})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get target staff")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}