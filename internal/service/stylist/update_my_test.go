package stylist

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestUpdateMyStylistService_UpdateMyStylist_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateMyStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	goodAtShapes := []string{"橢圓形", "方形"}
	isIntrovert := true

	req := stylist.UpdateMyStylistRequest{
		StylistName:  &stylistName,
		GoodAtShapes: &goodAtShapes,
		IsIntrovert:  &isIntrovert,
	}

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update
	expectedResponse := &stylist.UpdateMyStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Jane Updated",
		GoodAtShapes: []string{"橢圓形", "方形"},
		GoodAtColors: []string{"裸色系"},
		GoodAtStyles: []string{"簡約"},
		IsIntrovert:  true,
	}
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(expectedResponse, nil)

	response, err := service.UpdateMyStylist(ctx, req, staffUserID)

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

func TestUpdateMyStylistService_UpdateMyStylist_StylistNotCreated(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateMyStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateMyStylistRequest{
		StylistName: &stylistName,
	}

	// Mock stylist existence check - does not exist
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(false, nil)

	response, err := service.UpdateMyStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotCreated, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateMyStylistService_UpdateMyStylist_UpdateNoRows(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateMyStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateMyStylistRequest{
		StylistName: &stylistName,
	}

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update - no rows returned
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(nil, fmt.Errorf("no rows returned from update"))

	response, err := service.UpdateMyStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateMyStylistService_UpdateMyStylist_DatabaseError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateMyStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	stylistName := "Jane Updated"
	req := stylist.UpdateMyStylistRequest{
		StylistName: &stylistName,
	}

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update - database error
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(nil, errors.New("database error"))

	response, err := service.UpdateMyStylist(ctx, req, staffUserID)

	assert.Error(t, err)
	assert.Nil(t, response)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestUpdateMyStylistService_UpdateMyStylist_PartialUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockQuerier := mocks.NewMockQuerier()
	mockRepository := mocks.NewMockStylistRepository()
	service := NewUpdateMyStylistService(mockQuerier, mockRepository)

	ctx := context.Background()
	staffUserID := int64(12345)
	isIntrovert := true

	req := stylist.UpdateMyStylistRequest{
		IsIntrovert: &isIntrovert, // Only update this field
	}

	// Mock stylist existence check
	mockQuerier.On("CheckStylistExistsByStaffUserID", ctx, pgtype.Int8{Int64: staffUserID, Valid: true}).Return(true, nil)

	// Mock repository update
	expectedResponse := &stylist.UpdateMyStylistResponse{
		ID:           "18000000001",
		StaffUserID:  "12345",
		StylistName:  "Original Name",
		GoodAtShapes: []string{"方形"},
		GoodAtColors: []string{"裸色系"},
		GoodAtStyles: []string{"簡約"},
		IsIntrovert:  true,
	}
	mockRepository.On("UpdateStylist", ctx, staffUserID, req).Return(expectedResponse, nil)

	response, err := service.UpdateMyStylist(ctx, req, staffUserID)

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
