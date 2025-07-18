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
)

func TestUpdateTimeSlotService_UpdateTimeSlot_AllFieldsEmpty(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.UpdateTimeSlotRequest{
		// All fields are nil
	}

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValAllFieldsEmpty, serviceErr.Code)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_OnlyStartTimeProvided(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	startTime := "09:00"
	request := schedule.UpdateTimeSlotRequest{
		StartTime: &startTime,
		// EndTime is nil
	}

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotCannotUpdateSeparately, serviceErr.Code)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_OnlyEndTimeProvided(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	endTime := "12:00"
	request := schedule.UpdateTimeSlotRequest{
		EndTime: &endTime,
		// StartTime is nil
	}

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotCannotUpdateSeparately, serviceErr.Code)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_OnlyIsAvailableProvided(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

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

	isAvailable := false
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
		// StartTime and EndTime are nil
	}

	expectedResponse := &schedule.UpdateTimeSlotResponse{
		ID:          "5000000001",
		ScheduleID:  "4000000001",
		StartTime:   "09:00",
		EndTime:     "12:00",
		IsAvailable: false,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockTimeSlotRepo.On("UpdateTimeSlot", ctx, timeSlotID, request).Return(expectedResponse, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, "4000000001", response.ScheduleID)
	assert.Equal(t, false, response.IsAvailable)

	mockQuerier.AssertExpectations(t)
	mockTimeSlotRepo.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_InvalidScheduleID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	response, err := service.UpdateTimeSlot(ctx, "invalid", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_InvalidTimeSlotID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "invalid", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_TimeSlotNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	timeSlotID := int64(5000000001)

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(dbgen.TimeSlot{}, assert.AnError)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_TimeSlotNotBelongToSchedule(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

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

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_TimeSlotBooked(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

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

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotAlreadyBookedDoNotUpdate, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_StylistPermissionDenied(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		IsAvailable: pgtype.Bool{Bool: true, Valid: true}, // Available to pass booking check
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
		Role:   staff.RoleStylist, // Stylist role - can only modify their own schedules
	}

	isAvailable := true
	request := schedule.UpdateTimeSlotRequest{
		IsAvailable: &isAvailable,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_InvalidStartTimeFormat(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		StartTime:   pgtype.Time{Microseconds: 9 * 3600 * 1000000, Valid: true},
		EndTime:     pgtype.Time{Microseconds: 12 * 3600 * 1000000, Valid: true},
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

	invalidStartTime := "invalid"
	endTime := "12:00"
	request := schedule.UpdateTimeSlotRequest{
		StartTime: &invalidStartTime,
		EndTime:   &endTime,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_EndTimeBeforeStartTime(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		StartTime:   pgtype.Time{Microseconds: 9 * 3600 * 1000000, Valid: true},
		EndTime:     pgtype.Time{Microseconds: 12 * 3600 * 1000000, Valid: true},
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

	startTime := "15:00"
	endTime := "12:00" // End time before start time
	request := schedule.UpdateTimeSlotRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotService_UpdateTimeSlot_TimeSlotOverlap(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := mocks.NewMockTimeSlotRepository()
	service := NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)

	ctx := context.Background()
	scheduleID := int64(4000000001)
	timeSlotID := int64(5000000001)
	stylistID := int64(12345)
	storeID := int64(67890)

	mockTimeSlot := dbgen.TimeSlot{
		ID:          timeSlotID,
		ScheduleID:  scheduleID,
		StartTime:   pgtype.Time{Microseconds: 9 * 3600 * 1000000, Valid: true},
		EndTime:     pgtype.Time{Microseconds: 12 * 3600 * 1000000, Valid: true},
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

	startTime := "10:00"
	endTime := "14:00"
	request := schedule.UpdateTimeSlotRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	mockQuerier.On("GetTimeSlotByID", ctx, timeSlotID).Return(mockTimeSlot, nil)
	mockQuerier.On("GetScheduleByID", ctx, scheduleID).Return(mockSchedule, nil)
	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("CheckTimeSlotOverlapExcluding", ctx, dbgen.CheckTimeSlotOverlapExcludingParams{
		ScheduleID: scheduleID,
		ID:         timeSlotID,
		StartTime:  pgtype.Time{Microseconds: 10 * 3600 * 1000000, Valid: true},
		EndTime:    pgtype.Time{Microseconds: 14 * 3600 * 1000000, Valid: true},
	}).Return(true, nil) // Has overlap

	response, err := service.UpdateTimeSlot(ctx, "4000000001", "5000000001", request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleTimeConflict, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

// Note: Success test with database update is handled by integration tests
// The update logic is straightforward and tested implicitly in the actual application
