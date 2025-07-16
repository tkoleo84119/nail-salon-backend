package stylist

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestUpdateStylistService_UpdateStylist_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	goodAtShapes := []string{"橢圓形", "方形"}
	isIntrovert := true

	req := stylist.UpdateStylistRequest{
		StylistName:  &stylistName,
		GoodAtShapes: &goodAtShapes,
		IsIntrovert:  &isIntrovert,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update
	expectedResponse := &stylist.UpdateStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Jane Updated",
		GoodAtShapes: []string{"橢圓形", "方形"},
		GoodAtColors: []string{"裸色系"},
		GoodAtStyles: []string{"簡約"},
		IsIntrovert:  true,
	}
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(expectedResponse, nil)

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "18000000001", response.ID)
	assert.Equal(t, "12345", response.StaffUserID)
	assert.Equal(t, "Jane Updated", response.StylistName)
	assert.Equal(t, []string{"橢圓形", "方形"}, response.GoodAtShapes)
	assert.Equal(t, []string{"裸色系"}, response.GoodAtColors)
	assert.Equal(t, []string{"簡約"}, response.GoodAtStyles)
	assert.Equal(t, true, response.IsIntrovert)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_SuperAdminForbidden(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock staff user lookup - SUPER_ADMIN role
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleSuperAdmin,
	}, nil)

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_NoFieldsToUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	req := stylist.UpdateStylistRequest{
		// No fields to update
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValAllFieldsEmpty, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_StylistNotCreated(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check - does not exist
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(false, nil)

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotCreated, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_StaffUserNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock staff user lookup - not found
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{}, errors.New("staff user not found"))

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_UpdateNoRows(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update - no rows returned
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(nil, fmt.Errorf("no rows returned from update"))

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateStylistRequest{
		StylistName: &stylistName,
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update - database error
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(nil, errors.New("database error"))

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateStylistService_UpdateStylist_PartialUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	isIntrovert := true

	req := stylist.UpdateStylistRequest{
		IsIntrovert: &isIntrovert, // Only update this field
	}

	// Mock staff user lookup
	mockQuerier.On("GetStaffUserByID", ctx, staffUserID).Return(dbgen.StaffUser{
		ID:   staffUserID,
		Role: staff.RoleAdmin,
	}, nil)

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update
	expectedResponse := &stylist.UpdateStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Original Name",
		GoodAtShapes: []string{"方形"},
		GoodAtColors: []string{"裸色系"},
		GoodAtStyles: []string{"簡約"},
		IsIntrovert:  true,
	}
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(expectedResponse, nil)

	response, err := service.UpdateStylist(ctx, req, staffUserID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "18000000001", response.ID)
	assert.Equal(t, "12345", response.StaffUserID)
	assert.Equal(t, "Original Name", response.StylistName)
	assert.Equal(t, []string{"方形"}, response.GoodAtShapes)
	assert.Equal(t, []string{"裸色系"}, response.GoodAtColors)
	assert.Equal(t, []string{"簡約"}, response.GoodAtStyles)
	assert.Equal(t, true, response.IsIntrovert)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}