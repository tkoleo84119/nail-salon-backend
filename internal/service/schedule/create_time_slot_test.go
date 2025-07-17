package schedule

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

func init() {
	// Initialize snowflake for testing
	utils.InitSnowflake(1)
}

func TestCreateTimeSlotService_CreateTimeSlot_InvalidScheduleID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	response, err := service.CreateTimeSlot(ctx, "invalid", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotService_CreateTimeSlot_InvalidStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "invalid",
		EndTime:   "12:00",
	}

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotService_CreateTimeSlot_InvalidEndTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "invalid",
	}

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotService_CreateTimeSlot_EndTimeBeforeStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "12:00",
		EndTime:   "09:00", // End time before start time
	}

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateTimeSlotService_CreateTimeSlot_ScheduleNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(dbgen.Schedule{}, assert.AnError)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_StylistNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)

	mockSchedule := dbgen.Schedule{
		ID:        scheduleID,
		StylistID: stylistID,
		StoreID:   67890,
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(dbgen.Stylist{}, assert.AnError)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_StylistPermissionDenied(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)

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
		Role:   staff.RoleStylist, // Stylist role - can only modify their own schedules
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_StoreNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

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

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(dbgen.GetStoreByIDRow{}, assert.AnError)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_StoreNotActive(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

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
		IsActive: pgtype.Bool{Bool: false, Valid: true}, // Not active
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotActive, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_NoStoreAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

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
			{ID: "99999", Name: "Other Store"}, // Different store
		},
	}

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateTimeSlotService_CreateTimeSlot_TimeSlotOverlap(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewCreateTimeSlotService(mockQuerier)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

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

	request := schedule.CreateTimeSlotRequest{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("CheckTimeSlotOverlap", ctx, dbgen.CheckTimeSlotOverlapParams{
		ScheduleID: scheduleID,
		StartTime:  pgtype.Time{Microseconds: 9 * 3600 * 1000000, Valid: true},  // 09:00
		EndTime:    pgtype.Time{Microseconds: 12 * 3600 * 1000000, Valid: true}, // 12:00
	}).Return(true, nil) // Has overlap

	response, err := service.CreateTimeSlot(ctx, "4000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

// Note: Success test with database creation is handled by integration tests
// The creation logic is straightforward and tested implicitly in the actual application
