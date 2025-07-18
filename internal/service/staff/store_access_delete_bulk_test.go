package staff

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestDeleteStoreAccessBulkService_DeleteStoreAccessBulk_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessBulkService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessBulkRequest{
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
	response, err := service.DeleteStoreAccessBulk(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

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

func TestDeleteStoreAccessBulkService_DeleteStoreAccessBulk_CannotModifySelf(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessBulkService(mockQuerier)

	// Test data - same ID for target and creator
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(123456789) // Same as target
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessBulkRequest{
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
	response, err := service.DeleteStoreAccessBulk(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_NOT_UPDATE_SELF")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessBulkService_DeleteStoreAccessBulk_CannotModifySuperAdmin(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessBulkService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessBulkRequest{
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
	response, err := service.DeleteStoreAccessBulk(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessBulkService_DeleteStoreAccessBulk_StaffNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessBulkService(mockQuerier)

	// Test data
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2, 3}

	req := staff.DeleteStoreAccessBulkRequest{
		StoreIDs: []string{"1", "2"},
	}

	// Set up mock expectations - staff not found
	mockQuerier.On("GetStaffUserByID", mock.Anything, targetStaffID).Return(dbgen.StaffUser{}, sql.ErrNoRows)

	// Call service
	response, err := service.DeleteStoreAccessBulk(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "USER_STAFF_NOT_FOUND")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteStoreAccessBulkService_DeleteStoreAccessBulk_PermissionDenied(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteStoreAccessBulkService(mockQuerier)

	// Test data - creator doesn't have access to store being deleted
	targetID := "123456789"
	targetStaffID := int64(123456789)
	creatorID := int64(987654321)
	creatorRole := staff.RoleAdmin
	creatorStoreIDs := []int64{1, 2} // Only has access to store 1 and 2

	req := staff.DeleteStoreAccessBulkRequest{
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
	response, err := service.DeleteStoreAccessBulk(context.Background(), targetID, req, creatorID, creatorRole, creatorStoreIDs)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_PERMISSION_DENIED")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}
