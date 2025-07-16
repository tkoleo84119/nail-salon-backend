package staff

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestUpdateStaffMeService_UpdateStaffMe_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "new-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:       staffUserID,
		Username: "staff_amy",
		Email:    "old-email@example.com",
		Role:     staff.RoleAdmin,
	}, nil)

	// Mock email uniqueness check
	mockQuerier.On("CheckEmailUniqueForUpdate", ctx, dbgen.CheckEmailUniqueForUpdateParams{
		Email: email,
		ID:    staffUserID,
	}).Return(false, nil)

	// Mock repository update
	expectedResponse := &staff.UpdateStaffMeResponse{
		ID:       "12345",
		Username: "staff_amy",
		Email:    "new-email@example.com",
		Role:     staff.RoleAdmin,
	}
	mockRepository.On("UpdateStaffMe", ctx, staffUserID, req).Return(expectedResponse, nil)

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "12345", response.ID)
	assert.Equal(t, "staff_amy", response.Username)
	assert.Equal(t, "new-email@example.com", response.Email)
	assert.Equal(t, staff.RoleAdmin, response.Role)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateStaffMeService_UpdateStaffMe_NoFieldsToUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := staff.UpdateStaffMeRequest{
		// No fields to update
	}

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValAllFieldsEmpty, serviceErr.Code)
}

func TestUpdateStaffMeService_UpdateStaffMe_StaffUserNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "new-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup - not found
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{}, errors.New("staff user not found"))

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffMeService_UpdateStaffMe_EmailAlreadyExists(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "existing-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:       staffUserID,
		Username: "staff_amy",
		Email:    "old-email@example.com",
		Role:     staff.RoleAdmin,
	}, nil)

	// Mock email uniqueness check - email already exists
	mockQuerier.On("CheckEmailUniqueForUpdate", ctx, dbgen.CheckEmailUniqueForUpdateParams{
		Email: email,
		ID:    staffUserID,
	}).Return(true, nil)

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserEmailExists, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffMeService_UpdateStaffMe_EmailUniquenessCheckError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "new-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:       staffUserID,
		Username: "staff_amy",
		Email:    "old-email@example.com",
		Role:     staff.RoleAdmin,
	}, nil)

	// Mock email uniqueness check - database error
	mockQuerier.On("CheckEmailUniqueForUpdate", ctx, dbgen.CheckEmailUniqueForUpdateParams{
		Email: email,
		ID:    staffUserID,
	}).Return(false, errors.New("database error"))

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStaffMeService_UpdateStaffMe_UpdateNoRows(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "new-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:       staffUserID,
		Username: "staff_amy",
		Email:    "old-email@example.com",
		Role:     staff.RoleAdmin,
	}, nil)

	// Mock email uniqueness check
	mockQuerier.On("CheckEmailUniqueForUpdate", ctx, dbgen.CheckEmailUniqueForUpdateParams{
		Email: email,
		ID:    staffUserID,
	}).Return(false, nil)

	// Mock repository update - no rows returned
	mockRepository.On("UpdateStaffMe", ctx, staffUserID, req).Return(nil, errors.New("no rows returned from update"))

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateStaffMeService_UpdateStaffMe_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStaffUserRepository()
	service := NewUpdateStaffMeService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	email := "new-email@example.com"

	req := staff.UpdateStaffMeRequest{
		Email: &email,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:       staffUserID,
		Username: "staff_amy",
		Email:    "old-email@example.com",
		Role:     staff.RoleAdmin,
	}, nil)

	// Mock email uniqueness check
	mockQuerier.On("CheckEmailUniqueForUpdate", ctx, dbgen.CheckEmailUniqueForUpdateParams{
		Email: email,
		ID:    staffUserID,
	}).Return(false, nil)

	// Mock repository update - database error
	mockRepository.On("UpdateStaffMe", ctx, staffUserID, req).Return(nil, errors.New("database error"))

	response, err := service.UpdateStaffMe(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}