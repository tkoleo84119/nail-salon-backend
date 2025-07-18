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

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock no existing items (empty slice)
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateID", ctx, int64(6000000011)).Return([]dbgen.TimeSlotTemplateItem{}, nil)

	// Mock successful item creation
	mockQuerier.On("CreateTimeSlotTemplateItem", ctx, mock.AnythingOfType("dbgen.CreateTimeSlotTemplateItemParams")).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
	}, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "6000000011", response.TemplateID)
	assert.Equal(t, "09:00", response.StartTime)
	assert.Equal(t, "12:00", response.EndTime)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_InvalidTemplateID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "invalid"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_TemplateNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Mock template not found
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{}, errors.New("not found"))

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_InvalidStartTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "invalid",
		EndTime:   "12:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_InvalidEndTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "09:00",
		EndTime:   "invalid",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_EndTimeBeforeStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "15:00",
		EndTime:   "12:00", // End before start
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_TimeConflict(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "10:00",
		EndTime:   "14:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock existing items with overlapping time
	existingItems := []dbgen.TimeSlotTemplateItem{
		{
			ID:         6100000001,
			TemplateID: 6000000011,
			StartTime:  pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:    pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
	}
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateID", ctx, int64(6000000011)).Return(existingItems, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_AdjacentTimeSlots(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "12:00",
		EndTime:   "15:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock existing items with adjacent time (not overlapping)
	existingItems := []dbgen.TimeSlotTemplateItem{
		{
			ID:         6100000001,
			TemplateID: 6000000011,
			StartTime:  pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:    pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
	}
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateID", ctx, int64(6000000011)).Return(existingItems, nil)

	// Mock successful item creation
	mockQuerier.On("CreateTimeSlotTemplateItem", ctx, mock.AnythingOfType("dbgen.CreateTimeSlotTemplateItemParams")).Return(dbgen.TimeSlotTemplateItem{
		ID:         6100000003,
		TemplateID: 6000000011,
	}, nil)

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "6000000011", response.TemplateID)
	assert.Equal(t, "12:00", response.StartTime)
	assert.Equal(t, "15:00", response.EndTime)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_CreateTimeSlotTemplateItem_GetExistingItemsError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateItemService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := schedule.CreateTimeSlotTemplateItemRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Test Template",
	}, nil)

	// Mock database error when getting existing items
	mockQuerier.On("GetTimeSlotTemplateItemsByTemplateID", ctx, int64(6000000011)).Return(nil, errors.New("database error"))

	// Call service
	response, err := service.CreateTimeSlotTemplateItem(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotTemplateItemService_checkTimeConflicts_Success(t *testing.T) {
	service := &CreateTimeSlotTemplateItemService{}

	startTime, _ := time.Parse("15:04", "13:00")
	endTime, _ := time.Parse("15:04", "17:00")

	// Existing items that don't conflict
	existingItems := []dbgen.TimeSlotTemplateItem{
		{
			StartTime: pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:   pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
		{
			StartTime: pgtype.Time{Microseconds: int64(18*3600) * 1000000, Valid: true}, // 18:00
			EndTime:   pgtype.Time{Microseconds: int64(21*3600) * 1000000, Valid: true}, // 21:00
		},
	}

	err := service.checkTimeConflicts(startTime, endTime, existingItems)
	assert.NoError(t, err)
}

func TestCreateTimeSlotTemplateItemService_checkTimeConflicts_Overlap(t *testing.T) {
	service := &CreateTimeSlotTemplateItemService{}

	startTime, _ := time.Parse("15:04", "10:00")
	endTime, _ := time.Parse("15:04", "14:00")

	// Existing items that conflict
	existingItems := []dbgen.TimeSlotTemplateItem{
		{
			StartTime: pgtype.Time{Microseconds: int64(9*3600) * 1000000, Valid: true},  // 09:00
			EndTime:   pgtype.Time{Microseconds: int64(12*3600) * 1000000, Valid: true}, // 12:00
		},
	}

	err := service.checkTimeConflicts(startTime, endTime, existingItems)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)
}
