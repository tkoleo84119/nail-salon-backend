package store

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

// MockStoreRepository implements the StoreRepositoryInterface for testing
type MockStoreRepository struct {
	mock.Mock
}

func (m *MockStoreRepository) UpdateStore(ctx context.Context, storeID int64, req store.UpdateStoreRequest) (*store.UpdateStoreResponse, error) {
	args := m.Called(ctx, storeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.UpdateStoreResponse), args.Error(1)
}

func TestUpdateStoreService_UpdateStore_InvalidStoreID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name: stringPtr("Updated Store"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	response, err := service.UpdateStore(ctx, "invalid", req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_NoUpdates(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{} // No fields to update
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_InsufficientPermission(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		shouldError bool
	}{
		{
			name:        "Manager_cannot_update_store",
			role:        staff.RoleManager,
			shouldError: true,
		},
		{
			name:        "Stylist_cannot_update_store",
			role:        staff.RoleStylist,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuerier := mocks.NewMockQuerier()
			mockStoreRepo := &MockStoreRepository{}
			service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

			req := store.UpdateStoreRequest{
				Name: stringPtr("Updated Store"),
			}
			staffContext := common.StaffContext{
				UserID: "33333",
				Role:   tt.role,
			}

			response, err := service.UpdateStore(context.Background(), "8000000001", req, staffContext)

			assert.Nil(t, response)
			assert.Error(t, err)
			serviceErr, ok := err.(*errorCodes.ServiceError)
			assert.True(t, ok)
			assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

			mockQuerier.AssertExpectations(t)
			mockStoreRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateStoreService_UpdateStore_StoreNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name: stringPtr("Updated Store"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock store not found
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(dbgen.Store{}, errors.New("store not found"))

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_StoreInactive_CanUpdate(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name: stringPtr("Updated Store"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock inactive store - should still be updatable
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Existing Store",
		IsActive: pgtype.Bool{Bool: false, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock name uniqueness check
	mockQuerier.On("CheckStoreNameExistsExcluding", ctx, dbgen.CheckStoreNameExistsExcludingParams{
		Name: "Updated Store",
		ID:   8000000001,
	}).Return(false, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Updated Store",
		Address:  "",
		Phone:    "",
		IsActive: false, // Remains inactive
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000001), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_AdminNoStoreAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name: stringPtr("Updated Store"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "8000000002"},
		},
	}

	// Mock store exists
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Existing Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_NameAlreadyExists(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name: stringPtr("Duplicate Name"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock active store
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Existing Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock name already exists
	mockQuerier.On("CheckStoreNameExistsExcluding", ctx, dbgen.CheckStoreNameExistsExcludingParams{
		Name: "Duplicate Name",
		ID:   8000000001,
	}).Return(true, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_Success_SuperAdmin(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name:    stringPtr("Updated Store"),
		Address: stringPtr("Updated Address"),
		Phone:   stringPtr("02-87654321"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock active store
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Existing Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock name uniqueness check
	mockQuerier.On("CheckStoreNameExistsExcluding", ctx, dbgen.CheckStoreNameExistsExcludingParams{
		Name: "Updated Store",
		ID:   8000000001,
	}).Return(false, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Updated Store",
		Address:  "Updated Address",
		Phone:    "02-87654321",
		IsActive: true,
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000001), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_Success_AdminWithAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Address: stringPtr("Updated Address"),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "8000000001"},
		},
	}

	// Mock active store
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Existing Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Existing Store",
		Address:  "Updated Address",
		Phone:    "",
		IsActive: true,
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000001), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestUpdateStoreService_UpdateStore_Success_UpdateIsActive(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		IsActive: boolPtr(false), // Deactivating store
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock active store
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Store to Deactivate",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Store to Deactivate",
		Address:  "",
		Phone:    "",
		IsActive: false,
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000001), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)
	assert.False(t, response.IsActive) // Verify store was deactivated

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_Success_ReactivateStore(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		IsActive: boolPtr(true), // Reactivating store
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock inactive store (we can update inactive stores to reactivate them)
	storeDetail := dbgen.Store{
		ID:       8000000002,
		Name:     "Store to Reactivate",
		IsActive: pgtype.Bool{Bool: false, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000002)).Return(storeDetail, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000002",
		Name:     "Store to Reactivate",
		Address:  "",
		Phone:    "",
		IsActive: true,
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000002), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000002", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)
	assert.True(t, response.IsActive) // Verify store was reactivated

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}

func TestUpdateStoreService_UpdateStore_Success_MultipleFields(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockStoreRepo := &MockStoreRepository{}
	service := NewUpdateStoreService(mockQuerier, mockStoreRepo)

	ctx := context.Background()
	req := store.UpdateStoreRequest{
		Name:     stringPtr("Updated Store Name"),
		Address:  stringPtr("Updated Address"),
		Phone:    stringPtr("02-98765432"),
		IsActive: boolPtr(false),
	}
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	// Mock active store
	storeDetail := dbgen.Store{
		ID:       8000000001,
		Name:     "Original Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}
	mockQuerier.On("GetStoreDetailByID", ctx, int64(8000000001)).Return(storeDetail, nil)

	// Mock name uniqueness check
	mockQuerier.On("CheckStoreNameExistsExcluding", ctx, dbgen.CheckStoreNameExistsExcludingParams{
		Name: "Updated Store Name",
		ID:   8000000001,
	}).Return(false, nil)

	// Mock successful update
	expectedResponse := &store.UpdateStoreResponse{
		ID:       "8000000001",
		Name:     "Updated Store Name",
		Address:  "Updated Address",
		Phone:    "02-98765432",
		IsActive: false,
	}
	mockStoreRepo.On("UpdateStore", ctx, int64(8000000001), req).Return(expectedResponse, nil)

	response, err := service.UpdateStore(ctx, "8000000001", req, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, "Updated Store Name", response.Name)
	assert.Equal(t, "Updated Address", response.Address)
	assert.Equal(t, "02-98765432", response.Phone)
	assert.False(t, response.IsActive)

	mockQuerier.AssertExpectations(t)
	mockStoreRepo.AssertExpectations(t)
}
