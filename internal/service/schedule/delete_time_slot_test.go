package schedule

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestDeleteTimeSlotService_DeleteTimeSlot_InvalidScheduleID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	response, err := service.DeleteTimeSlot(ctx, "invalid", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_InvalidTimeSlotID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "invalid", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_TimeSlotNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	timeSlotID := int64(5000000001)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(dbgen.TimeSlot{}, assert.AnError)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_TimeSlotNotBelongToSchedule(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	timeSlotID := int64(5000000001)

	mockTimeSlot := dbgen.TimeSlot{
		ID:         timeSlotID,
		ScheduleID: int64(9999999999), // Different schedule ID
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_TimeSlotBooked(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: false, Valid: true}, // Time slot is not available (booked)
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotAlreadyBookedDoNotDelete, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_ScheduleNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(dbgen.Schedule{}, assert.AnError)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_StylistNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   67890,
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(dbgen.Stylist{}, assert.AnError)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_StylistPermissionDenied(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   67890,
	}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 22222, Valid: true}, // Different staff user
	}

	staffContext := common.StaffContext{
		UserID: "11111",           // Current user
		Role:   staff.RoleStylist, // Stylist role - can only delete their own schedules
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_StoreNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   storeID,
	}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(dbgen.GetStoreByIDRow{}, assert.AnError)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_StoreNotActive(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   storeID,
	}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: false, Valid: true}, // Store is not active
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotActive, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_NoStoreAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   storeID,
	}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "99999", Name: "Other Store"}, // No access to store 67890
		},
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotService_DeleteTimeSlot_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true},
	}

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   storeID,
	}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("DeleteTimeSlotByID", ctx, timeSlotID).Return(nil)

	response, err := service.DeleteTimeSlot(ctx, "4000000001", "5000000001", staffContext)

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, []string{"5000000001"}, response.Deleted)

	mockQuerier.AssertExpectations(t)
}