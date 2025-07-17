package schedule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestCreateTimeSlotTemplateService_CreateTimeSlotTemplate_InvalidTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		TimeSlots: []schedule.TimeSlotItem{
			{StartTime: "invalid", EndTime: "12:00"},
		},
	}

	response, err := service.CreateTimeSlotTemplate(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotTemplateService_CreateTimeSlotTemplate_EndTimeBeforeStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		TimeSlots: []schedule.TimeSlotItem{
			{StartTime: "15:00", EndTime: "12:00"}, // End before start
		},
	}

	response, err := service.CreateTimeSlotTemplate(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotTemplateService_CreateTimeSlotTemplate_OverlappingTimeSlots(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		TimeSlots: []schedule.TimeSlotItem{
			{StartTime: "09:00", EndTime: "12:00"},
			{StartTime: "11:00", EndTime: "14:00"}, // Overlaps with first slot
		},
	}

	response, err := service.CreateTimeSlotTemplate(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)
}

func TestCreateTimeSlotTemplateService_CreateTimeSlotTemplate_EmptyTimeSlots(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotTemplateRequest{
		Name:      "Test Template",
		TimeSlots: []schedule.TimeSlotItem{}, // Empty time slots
	}

	response, err := service.CreateTimeSlotTemplate(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotTemplateService_CreateTimeSlotTemplate_InvalidStaffUserID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "invalid", // Invalid ID
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotTemplateRequest{
		Name: "Test Template",
		TimeSlots: []schedule.TimeSlotItem{
			{StartTime: "09:00", EndTime: "12:00"},
		},
	}

	response, err := service.CreateTimeSlotTemplate(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthStaffFailed, serviceErr.Code)
}

func TestCreateTimeSlotTemplateService_validateTimeSlots_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	timeSlots := []schedule.TimeSlotItem{
		{StartTime: "09:00", EndTime: "12:00"},
		{StartTime: "13:00", EndTime: "17:00"},
	}

	err := service.validateTimeSlots(timeSlots)

	assert.NoError(t, err)
}

func TestCreateTimeSlotTemplateService_validateTimeSlots_AdjacentSlots(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotTemplateService(mockQuerier, nil)

	timeSlots := []schedule.TimeSlotItem{
		{StartTime: "09:00", EndTime: "12:00"},
		{StartTime: "12:00", EndTime: "15:00"}, // Adjacent (not overlapping)
	}

	err := service.validateTimeSlots(timeSlots)

	assert.NoError(t, err)
}
