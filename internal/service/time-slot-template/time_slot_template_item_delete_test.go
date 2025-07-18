package timeSlotTemplate

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item exists and belongs to template
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
	}, nil)

	// Mock successful deletion
	mockQuerier.On("DeleteTimeSlotTemplateItem", ctx, int64(6100000003)).Return(nil)

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []string{"6100000003"}, response.Deleted)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_InvalidTemplateID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "invalid"
	itemID := "6100000003"

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_InvalidItemID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "invalid"

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_TemplateNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"

	// Mock template not found
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{}, errors.New("not found"))

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_ItemNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item not found
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{}, errors.New("not found"))

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateItemNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_ItemNotInTemplate(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item exists but belongs to different template
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000012, // Different template
	}, nil)

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateItemNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateItemService_DeleteTimeSlotTemplateItem_DeleteError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item exists and belongs to template
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
	}, nil)

	// Mock deletion error
	mockQuerier.On("DeleteTimeSlotTemplateItem", ctx, int64(6100000003)).Return(errors.New("database error"))

	// Call service
	response, err := service.DeleteTimeSlotTemplateItem(ctx, templateID, itemID, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}