package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item exists and belongs to template
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
		StartTime:  pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
		EndTime:    pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
	}, nil)

	// Mock no other conflicting items
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateIDExcluding", ctx, mock.AnythingOfType("dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams")).Return([]dbgen.TimeSlotTemplateItem{}, nil)

	// Mock successful update
	mockQuerier.On("UpdateTimeSlotTemplateItem", ctx, mock.AnythingOfType("dbgen.UpdateTimeSlotTemplateItemParams")).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
	}, nil)

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "6100000003", response.ID)
	assert.Equal(t, "6000000011", response.TemplateID)
	assert.Equal(t, "14:00", response.StartTime)
	assert.Equal(t, "18:00", response.EndTime)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_InvalidTemplateID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "invalid"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_InvalidItemID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "invalid"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_TemplateNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Mock template not found
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{}, errors.New("not found"))

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_ItemNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock item not found
	mockQuerier.On("GetTimeSlotTemplateItemByID", ctx, int64(6100000003)).Return(dbgen.TimeSlotTemplateItem{}, errors.New("not found"))

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateItemNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_ItemNotInTemplate(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

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
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateItemNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_InvalidStartTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "invalid",
		EndTime:   "18:00",
	}

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

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_InvalidEndTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "invalid",
	}

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

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_EndTimeBeforeStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "18:00",
		EndTime:   "14:00", // End before start
	}

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

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_TimeConflictWithOtherItems(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "14:00",
	}

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

	// Mock other items with overlapping time
	otherItems := []dbgen.TimeSlotTemplateItem{
		{
			ID:         6100000001,
			TemplateID: 6000000011,
			StartTime:  pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:    pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
	}
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateIDExcluding", ctx, mock.AnythingOfType("dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams")).Return(otherItems, nil)

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_GetOtherItemsError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

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

	// Mock database error when getting other items
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateIDExcluding", ctx, mock.AnythingOfType("dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams")).Return(nil, errors.New("database error"))

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_UpdateTimeSlotTemplateItem_UpdateError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewUpdateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	itemID := "6100000003"
	request := schedule.UpdateTimeSlotTemplateItemRequest{
		StartTime: "14:00",
		EndTime:   "18:00",
	}

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

	// Mock no conflicting items
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateIDExcluding", ctx, mock.AnythingOfType("dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams")).Return([]dbgen.TimeSlotTemplateItem{}, nil)

	// Mock update error
	mockQuerier.On("UpdateTimeSlotTemplateItem", ctx, mock.AnythingOfType("dbgen.UpdateTimeSlotTemplateItemParams")).Return(dbgen.TimeSlotTemplateItem{}, errors.New("update error"))

	// Call service
	response, err := service.UpdateTimeSlotTemplateItem(ctx, templateID, itemID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateItemService_checkTimeConflicts_NoConflict(t *testing.T) {
	service := &UpdateTimeSlotTemplateItemService{}

	startTime, _ := time.Parse("15:04", "13:00")
	endTime, _ := time.Parse("15:04", "17:00")

	// Other items that don't conflict
	otherItems := []dbgen.TimeSlotTemplateItem{
		{
			StartTime: pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:   pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
		{
			StartTime: pgtype.Time{Microseconds: int64(18*3600) * 1000000, Valid: true}, // 18:00
			EndTime:   pgtype.Time{Microseconds: int64(21*3600) * 1000000, Valid: true}, // 21:00
		},
	}

	err := service.checkTimeConflicts(startTime, endTime, otherItems)
	assert.NoError(t, err)
}

func TestUpdateTimeSlotTemplateItemService_checkTimeConflicts_WithConflict(t *testing.T) {
	service := &UpdateTimeSlotTemplateItemService{}

	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "14:00")

	// Other items that conflict
	otherItems := []dbgen.TimeSlotTemplateItem{
		{
			StartTime: pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:   pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
	}

	err := service.checkTimeConflicts(startTime, endTime, otherItems)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)
}