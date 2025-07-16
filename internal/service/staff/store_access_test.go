package staff

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
	"github.com/jackc/pgx/v5/pgtype"
)


func TestCreateStoreAccessService_CreateStoreAccess_Success_NewlyCreated(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.NoError(t, err)
	assert.True(t, isNewlyCreated)
	assert.NotNil(t, response)
	assert.Equal(t, "123456789", response.StaffUserID)
	assert.Len(t, response.StoreList, 2)
	assert.Equal(t, "1", response.StoreList[0].ID)
	assert.Equal(t, "Store 1", response.StoreList[0].Name)
	assert.Equal(t, "2", response.StoreList[1].ID)
	assert.Equal(t, "Test Store", response.StoreList[1].Name)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_Success_AlreadyExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

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
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Call service with invalid ID
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "invalid_id", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid target staff ID")
}

func TestCreateStoreAccessService_CreateStoreAccess_StaffNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Mock expectations - staff not found
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, sql.ErrNoRows)

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STAFF_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_CannotUpdateSelf(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "1", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_UPDATE_SELF")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_CannotUpdateSuperAdmin(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_StoreNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_StoreNotActive(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STORE_NOT_ACTIVE")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_AdminNoPermission(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "3", // Admin doesn't have access to store 3
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
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestCreateStoreAccessService_CreateStoreAccess_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := &CreateStoreAccessService{queries: mockQuerier}

	// Test data
	req := staff.CreateStoreAccessRequest{
		StoreID: "2",
	}

	// Mock expectations - database error
	mockQuerier.On("GetStaffUserByID", mock.Anything, int64(123456789)).Return(dbgen.StaffUser{}, errors.New("database connection failed"))

	// Call service
	response, isNewlyCreated, err := service.CreateStoreAccess(context.Background(), "123456789", req, int64(1), staff.RoleSuperAdmin, []int64{1})

	// Assert results
	assert.Error(t, err)
	assert.False(t, isNewlyCreated)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get target staff")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessService_DeleteStoreAccess_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2"},
	}

	// Mock target staff user
	targetStaff := dbgen.StaffUser{
		ID:       targetStaffID,
		Username: "testuser",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Mock remaining store access after deletion
	remainingStoreAccess := []dbgen.GetStaffUserStoreAccessRow{
		{StoreID: 3, StoreName: "Store 3"},
	}

	// Set up mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(targetStaff, nil)
	mockQuerier.On("DeleteStaffUserStoreAccess", mock.Anything, mock.AnythingOfType("dbgen.DeleteStaffUserStoreAccessParams")).Return(nil)
	mockQuerier.On("GetStaffUserStoreAccess", mock.Anything, targetStaffID).Return(remainingStoreAccess, nil)

	// Call service
	response, err := service.DeleteStoreAccess(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "123456789", response.StaffUserID)
	assert.Len(t, response.StoreList, 1)
	assert.Equal(t, "3", response.StoreList[0].ID)
	assert.Equal(t, "Store 3", response.StoreList[0].Name)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessService_DeleteStoreAccess_CannotModifySelf(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessService(mockQuerier)

	// Test data - same ID for target and creator
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(123456789) // Same as target
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2"},
	}

	// Mock target staff user
	targetStaff := dbgen.StaffUser{
		ID:       targetStaffID,
		Username: "testuser",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Set up mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(targetStaff, nil)

	// Call service
	response, err := service.DeleteStoreAccess(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_UPDATE_SELF")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessService_DeleteStoreAccess_CannotModifySuperAdmin(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2"},
	}

	// Mock target staff user as SUPER_ADMIN
	targetStaff := dbgen.StaffUser{
		ID:       targetStaffID,
		Username: "superadmin",
		Role:     staff.RoleSuperAdmin,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Set up mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(targetStaff, nil)

	// Call service
	response, err := service.DeleteStoreAccess(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_CANNOT_MODIFY_SUPER_ADMIN")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessService_DeleteStoreAccess_StaffNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "2"},
	}

	// Set up mock expectations - staff not found
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(dbgen.StaffUser{}, sql.ErrNoRows)

	// Call service
	response, err := service.DeleteStoreAccess(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STAFF_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessService_DeleteStoreAccess_PermissionDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessService(mockQuerier)

	// Test data - creator doesn't have access to store being deleted
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2} // Only has access to store 1 and 2

	req := staff.DeleteStoreAccessRequest{
		StoreIDs: []string{"1", "3"}, // Trying to delete store 3 which creator doesn't have access to
	}

	// Mock target staff user
	targetStaff := dbgen.StaffUser{
		ID:       targetStaffID,
		Username: "testuser",
		Role:     staff.RoleStylist,
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Set up mock expectations
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(targetStaff, nil)

	// Call service
	response, err := service.DeleteStoreAccess(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}