package schedule

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

func TestCreateSchedulesBulkService_prepareBatchData_DirectTest(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	storeID := int64(67890)
	stylistID := int64(12345)

	note := "Test note"
	schedules := []schedule.ScheduleRequest{
		{
			WorkDate: "2023-12-01",
			Note:     &note,
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "09:00", EndTime: "10:00"},
				{StartTime: "14:00", EndTime: "15:00"},
			},
		},
	}

	// Test the prepareBatchData method directly
	scheduleRows, timeSlotRows, createdScheduleIDs, err := service.prepareBatchData(schedules, storeID, stylistID)

	assert.NoError(t, err)
	assert.Len(t, scheduleRows, 1)
	assert.Len(t, timeSlotRows, 2)
	assert.Len(t, createdScheduleIDs, 1)

	// Validate schedule row
	scheduleRow := scheduleRows[0]
	assert.Equal(t, storeID, scheduleRow.StoreID)
	assert.Equal(t, stylistID, scheduleRow.StylistID)
	assert.Equal(t, "Test note", scheduleRow.Note.String)
	assert.True(t, scheduleRow.Note.Valid)

	// Validate time slot rows
	assert.Equal(t, scheduleRow.ID, timeSlotRows[0].ScheduleID)
	assert.Equal(t, scheduleRow.ID, timeSlotRows[1].ScheduleID)
	assert.True(t, timeSlotRows[0].IsAvailable.Bool)
	assert.True(t, timeSlotRows[1].IsAvailable.Bool)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_InvalidStylistID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "invalid",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_InvalidStoreID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "invalid",
		Schedules: []schedule.ScheduleRequest{},
	}

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_StylistNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(dbgen.Stylist{}, assert.AnError)

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_StylistPermissionDenied(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	otherStaffUserID := int64(22222)

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: otherStaffUserID, Valid: true}, // Different staff user
	}

	staffContext := common.StaffContext{
		UserID: "11111",           // Current user
		Role:   staff.RoleStylist, // Stylist role - can only create for themselves
	}

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_StoreNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(dbgen.GetStoreByIDRow{}, assert.AnError)

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_StoreNotActive(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)

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

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotActive, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateSchedulesBulkService_CreateSchedulesBulk_NoStoreAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewCreateSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)

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

	request := schedule.CreateSchedulesBulkRequest{
		StylistID: "12345",
		StoreID:   "67890",
		Schedules: []schedule.ScheduleRequest{},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.CreateSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestCreateSchedulesBulkService_validateSchedules_DuplicateWorkDate(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	schedules := []schedule.ScheduleRequest{
		{
			WorkDate: "2023-12-01",
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "09:00", EndTime: "10:00"},
			},
		},
		{
			WorkDate: "2023-12-01", // Duplicate
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "14:00", EndTime: "15:00"},
			},
		},
	}

	err := service.validateSchedules(schedules)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_Success(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "09:00", EndTime: "10:00"},
		{StartTime: "14:00", EndTime: "15:00"},
		{StartTime: "16:00", EndTime: "17:00"},
	}

	err := service.validateTimeSlots(timeSlots)

	assert.NoError(t, err)
}

func TestCreateSchedulesBulkService_validateTimeSlots_EmptySlots(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{}

	err := service.validateTimeSlots(timeSlots)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_InvalidStartTime(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "invalid", EndTime: "10:00"},
	}

	err := service.validateTimeSlots(timeSlots)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_InvalidEndTime(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "09:00", EndTime: "invalid"},
	}

	err := service.validateTimeSlots(timeSlots)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_EndBeforeStart(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "10:00", EndTime: "09:00"}, // End before start
	}

	err := service.validateTimeSlots(timeSlots)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_OverlappingSlots(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "09:00", EndTime: "10:30"},
		{StartTime: "10:00", EndTime: "11:00"}, // Overlaps with first slot
	}

	err := service.validateTimeSlots(timeSlots)

	assert.Error(t, err)
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)
}

func TestCreateSchedulesBulkService_validateTimeSlots_AdjacentSlots(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	timeSlots := []schedule.TimeSlotRequest{
		{StartTime: "09:00", EndTime: "10:00"},
		{StartTime: "10:00", EndTime: "11:00"}, // Adjacent - should be valid
	}

	err := service.validateTimeSlots(timeSlots)

	assert.NoError(t, err)
}

func TestCreateSchedulesBulkService_prepareBatchData_Success(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	note := "Test note"
	schedules := []schedule.ScheduleRequest{
		{
			WorkDate: "2023-12-01",
			Note:     &note,
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "09:00", EndTime: "10:00"},
				{StartTime: "14:00", EndTime: "15:00"},
			},
		},
		{
			WorkDate: "2023-12-02",
			Note:     nil,
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "11:00", EndTime: "12:00"},
			},
		},
	}

	storeID := int64(67890)
	stylistID := int64(12345)

	scheduleRows, timeSlotRows, createdScheduleIDs, err := service.prepareBatchData(schedules, storeID, stylistID)

	assert.NoError(t, err)
	assert.Len(t, scheduleRows, 2)
	assert.Len(t, timeSlotRows, 3) // 2 + 1 time slots
	assert.Len(t, createdScheduleIDs, 2)

	// Check first schedule
	assert.Equal(t, storeID, scheduleRows[0].StoreID)
	assert.Equal(t, stylistID, scheduleRows[0].StylistID)
	assert.True(t, scheduleRows[0].Note.Valid)
	assert.Equal(t, "Test note", scheduleRows[0].Note.String)

	// Check second schedule (no note)
	assert.Equal(t, storeID, scheduleRows[1].StoreID)
	assert.Equal(t, stylistID, scheduleRows[1].StylistID)
	assert.False(t, scheduleRows[1].Note.Valid)

	// Check time slots
	for _, timeSlotRow := range timeSlotRows {
		assert.True(t, timeSlotRow.IsAvailable.Bool)
		assert.True(t, timeSlotRow.CreatedAt.Valid)
		assert.True(t, timeSlotRow.UpdatedAt.Valid)
	}
}

func TestCreateSchedulesBulkService_prepareBatchData_InvalidWorkDate(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	schedules := []schedule.ScheduleRequest{
		{
			WorkDate: "invalid-date",
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "09:00", EndTime: "10:00"},
			},
		},
	}

	storeID := int64(67890)
	stylistID := int64(12345)

	scheduleRows, timeSlotRows, createdScheduleIDs, err := service.prepareBatchData(schedules, storeID, stylistID)

	assert.Nil(t, scheduleRows)
	assert.Nil(t, timeSlotRows)
	assert.Nil(t, createdScheduleIDs)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestCreateSchedulesBulkService_prepareBatchData_InvalidTimeSlot(t *testing.T) {
	service := &CreateSchedulesBulkService{}

	schedules := []schedule.ScheduleRequest{
		{
			WorkDate: "2023-12-01",
			TimeSlots: []schedule.TimeSlotRequest{
				{StartTime: "invalid", EndTime: "10:00"},
			},
		},
	}

	storeID := int64(67890)
	stylistID := int64(12345)

	scheduleRows, timeSlotRows, createdScheduleIDs, err := service.prepareBatchData(schedules, storeID, stylistID)

	assert.Nil(t, scheduleRows)
	assert.Nil(t, timeSlotRows)
	assert.Nil(t, createdScheduleIDs)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}
