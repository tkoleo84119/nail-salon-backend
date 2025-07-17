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

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_InvalidStylistID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "invalid",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_InvalidStoreID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "invalid",
		ScheduleIDs: []string{"4000000001"},
	}

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_InvalidScheduleIDs(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"invalid"},
	}

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_StylistNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(dbgen.Stylist{}, assert.AnError)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.StylistNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_StylistPermissionDenied(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	otherStaffUserID := int64(22222)

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: otherStaffUserID, Valid: true}, // Different staff user
	}

	staffContext := common.StaffContext{
		UserID: "11111", // Current user
		Role:   staff.RoleStylist, // Stylist role - can only delete their own schedules
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_StoreNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

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

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(dbgen.GetStoreByIDRow{}, assert.AnError)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_StoreNotActive(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

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

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.UserStoreNotActive, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_NoStoreAccess(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

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

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.AuthPermissionDenied, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_SchedulesNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)
	scheduleIDs := []int64{4000000001, 4000000002}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// Only 1 schedule found instead of 2
	mockSchedules := []dbgen.GetSchedulesWithTimeSlotsByIDsRow{
		{
			ID:        4000000001,
			StoreID:   storeID,
			StylistID: stylistID,
			IsAvailable: pgtype.Bool{Bool: true, Valid: true},
		},
		// Missing second schedule
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001", "4000000002"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("GetSchedulesWithTimeSlotsByIDs", ctx, scheduleIDs).Return(mockSchedules, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_SchedulesNotOwned(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)
	scheduleIDs := []int64{4000000001, 4000000002}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// All schedules found but one belongs to different stylist/store
	mockSchedules := []dbgen.GetSchedulesWithTimeSlotsByIDsRow{
		{
			ID:        4000000001,
			StoreID:   storeID,
			StylistID: stylistID,
			IsAvailable: pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:        4000000002,
			StoreID:   99999, // Different store
			StylistID: stylistID,
			IsAvailable: pgtype.Bool{Bool: true, Valid: true},
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001", "4000000002"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("GetSchedulesWithTimeSlotsByIDs", ctx, scheduleIDs).Return(mockSchedules, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteSchedulesBulkService_DeleteSchedulesBulk_SchedulesAlreadyBooked(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockDB := &pgxpool.Pool{}
	service := NewDeleteSchedulesBulkService(mockQuerier, mockDB)

	ctx := context.Background()
	stylistID := int64(12345)
	storeID := int64(67890)
	scheduleIDs := []int64{4000000001, 4000000002}

	mockStylist := dbgen.Stylist{
		ID:          stylistID,
		StaffUserID: pgtype.Int8{Int64: 11111, Valid: true},
	}

	mockStore := dbgen.GetStoreByIDRow{
		ID:       storeID,
		Name:     "Test Store",
		IsActive: pgtype.Bool{Bool: true, Valid: true},
	}

	// All schedules found but one is already booked
	mockSchedules := []dbgen.GetSchedulesWithTimeSlotsByIDsRow{
		{
			ID:        4000000001,
			StoreID:   storeID,
			StylistID: stylistID,
			IsAvailable: pgtype.Bool{Bool: true, Valid: true},
		},
		{
			ID:        4000000002,
			StoreID:   storeID,
			StylistID: stylistID,
			IsAvailable: pgtype.Bool{Bool: false, Valid: true}, // Already booked
		},
	}

	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
		StoreList: []common.Store{
			{ID: "67890", Name: "Test Store"},
		},
	}

	request := schedule.DeleteSchedulesBulkRequest{
		StylistID:   "12345",
		StoreID:     "67890",
		ScheduleIDs: []string{"4000000001", "4000000002"},
	}

	mockQuerier.On("GetStylistByID", ctx, stylistID).Return(mockStylist, nil)
	mockQuerier.On("GetStoreByID", ctx, storeID).Return(mockStore, nil)
	mockQuerier.On("GetSchedulesWithTimeSlotsByIDs", ctx, scheduleIDs).Return(mockSchedules, nil)

	response, err := service.DeleteSchedulesBulk(ctx, request, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)
	
	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ScheduleAlreadyBookedDoNotDelete, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

// Note: Success test with database transaction is handled by integration tests
// The transaction logic is straightforward and tested implicitly in the actual application