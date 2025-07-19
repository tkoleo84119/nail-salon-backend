package mocks

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

// MockQuerier implements the dbgen.Querier interface for testing
type MockQuerier struct {
	mock.Mock
}

// Ensure MockQuerier implements the interface
var _ dbgen.Querier = (*MockQuerier)(nil)

// NewMockQuerier creates a new instance of MockQuerier
func NewMockQuerier() *MockQuerier {
	return &MockQuerier{}
}

// Staff User related mock methods
func (m *MockQuerier) GetStaffUserByID(ctx context.Context, userID int64) (dbgen.StaffUser, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) GetStaffUserByUsername(ctx context.Context, username string) (dbgen.StaffUser, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) CreateStaffUser(ctx context.Context, arg dbgen.CreateStaffUserParams) (dbgen.CreateStaffUserRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.CreateStaffUserRow), args.Error(1)
}

func (m *MockQuerier) CheckStaffUserExists(ctx context.Context, arg dbgen.CheckStaffUserExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CheckEmailUniqueForUpdate(ctx context.Context, arg dbgen.CheckEmailUniqueForUpdateParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// Staff User Token related mock methods
func (m *MockQuerier) CreateStaffUserToken(ctx context.Context, arg dbgen.CreateStaffUserTokenParams) (dbgen.CreateStaffUserTokenRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.CreateStaffUserTokenRow), args.Error(1)
}

// Store related mock methods
func (m *MockQuerier) GetAllActiveStores(ctx context.Context) ([]dbgen.GetAllActiveStoresRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dbgen.GetAllActiveStoresRow), args.Error(1)
}

func (m *MockQuerier) GetStoreByID(ctx context.Context, id int64) (dbgen.GetStoreByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.GetStoreByIDRow), args.Error(1)
}

func (m *MockQuerier) GetStoresByIDs(ctx context.Context, storeIDs []int64) ([]dbgen.GetStoresByIDsRow, error) {
	args := m.Called(ctx, storeIDs)
	return args.Get(0).([]dbgen.GetStoresByIDsRow), args.Error(1)
}

func (m *MockQuerier) CheckStoresExistAndActive(ctx context.Context, storeIDs []int64) (dbgen.CheckStoresExistAndActiveRow, error) {
	args := m.Called(ctx, storeIDs)
	return args.Get(0).(dbgen.CheckStoresExistAndActiveRow), args.Error(1)
}

// Staff User Store Access related mock methods
func (m *MockQuerier) GetStaffUserStoreAccess(ctx context.Context, staffUserID int64) ([]dbgen.GetStaffUserStoreAccessRow, error) {
	args := m.Called(ctx, staffUserID)
	return args.Get(0).([]dbgen.GetStaffUserStoreAccessRow), args.Error(1)
}

func (m *MockQuerier) CreateStaffUserStoreAccess(ctx context.Context, arg dbgen.CreateStaffUserStoreAccessParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) BatchCreateStaffUserStoreAccess(ctx context.Context, arg []dbgen.BatchCreateStaffUserStoreAccessParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) CheckStoreAccessExists(ctx context.Context, arg dbgen.CheckStoreAccessExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) DeleteStaffUserStoreAccess(ctx context.Context, arg dbgen.DeleteStaffUserStoreAccessParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// Stylist related mock methods
func (m *MockQuerier) GetStylistByStaffUserID(ctx context.Context, staffUserID pgtype.Int8) (dbgen.Stylist, error) {
	args := m.Called(ctx, staffUserID)
	return args.Get(0).(dbgen.Stylist), args.Error(1)
}

func (m *MockQuerier) CheckStylistExistsByStaffUserID(ctx context.Context, staffUserID pgtype.Int8) (bool, error) {
	args := m.Called(ctx, staffUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CreateStylist(ctx context.Context, arg dbgen.CreateStylistParams) (dbgen.Stylist, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.Stylist), args.Error(1)
}

// Schedule related mock methods
func (m *MockQuerier) GetStylistByID(ctx context.Context, id int64) (dbgen.Stylist, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.Stylist), args.Error(1)
}

func (m *MockQuerier) CheckScheduleExists(ctx context.Context, arg dbgen.CheckScheduleExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CreateSchedule(ctx context.Context, arg dbgen.CreateScheduleParams) (dbgen.Schedule, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.Schedule), args.Error(1)
}

func (m *MockQuerier) CreateTimeSlot(ctx context.Context, arg dbgen.CreateTimeSlotParams) (dbgen.TimeSlot, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.TimeSlot), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotsByScheduleID(ctx context.Context, scheduleID int64) ([]dbgen.TimeSlot, error) {
	args := m.Called(ctx, scheduleID)
	return args.Get(0).([]dbgen.TimeSlot), args.Error(1)
}

func (m *MockQuerier) BatchCreateSchedules(ctx context.Context, arg []dbgen.BatchCreateSchedulesParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) BatchCreateTimeSlots(ctx context.Context, arg []dbgen.BatchCreateTimeSlotsParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetSchedulesByStoreAndStylist(ctx context.Context, arg dbgen.GetSchedulesByStoreAndStylistParams) ([]dbgen.Schedule, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]dbgen.Schedule), args.Error(1)
}

func (m *MockQuerier) GetSchedulesWithTimeSlotsByIDs(ctx context.Context, scheduleIDs []int64) ([]dbgen.GetSchedulesWithTimeSlotsByIDsRow, error) {
	args := m.Called(ctx, scheduleIDs)
	return args.Get(0).([]dbgen.GetSchedulesWithTimeSlotsByIDsRow), args.Error(1)
}

func (m *MockQuerier) DeleteSchedulesByIDs(ctx context.Context, scheduleIDs []int64) error {
	args := m.Called(ctx, scheduleIDs)
	return args.Error(0)
}

func (m *MockQuerier) DeleteTimeSlotsByScheduleIDs(ctx context.Context, scheduleIDs []int64) error {
	args := m.Called(ctx, scheduleIDs)
	return args.Error(0)
}

func (m *MockQuerier) GetScheduleByID(ctx context.Context, id int64) (dbgen.Schedule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.Schedule), args.Error(1)
}

func (m *MockQuerier) CheckTimeSlotOverlap(ctx context.Context, arg dbgen.CheckTimeSlotOverlapParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotByID(ctx context.Context, id int64) (dbgen.TimeSlot, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.TimeSlot), args.Error(1)
}

func (m *MockQuerier) CheckTimeSlotOverlapExcluding(ctx context.Context, arg dbgen.CheckTimeSlotOverlapExcludingParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) DeleteTimeSlotByID(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Time Slot Template related mock methods
func (m *MockQuerier) CreateTimeSlotTemplate(ctx context.Context, arg dbgen.CreateTimeSlotTemplateParams) (dbgen.TimeSlotTemplate, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.TimeSlotTemplate), args.Error(1)
}

func (m *MockQuerier) CreateTimeSlotTemplateItem(ctx context.Context, arg dbgen.CreateTimeSlotTemplateItemParams) (dbgen.TimeSlotTemplateItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.TimeSlotTemplateItem), args.Error(1)
}

func (m *MockQuerier) BatchCreateTimeSlotTemplateItems(ctx context.Context, arg []dbgen.BatchCreateTimeSlotTemplateItemsParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotTemplateByID(ctx context.Context, id int64) (dbgen.TimeSlotTemplate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.TimeSlotTemplate), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotTemplateItemsByTemplateID(ctx context.Context, templateID int64) ([]dbgen.TimeSlotTemplateItem, error) {
	args := m.Called(ctx, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dbgen.TimeSlotTemplateItem), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotTemplateItemByID(ctx context.Context, id int64) (dbgen.TimeSlotTemplateItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.TimeSlotTemplateItem), args.Error(1)
}

func (m *MockQuerier) UpdateTimeSlotTemplateItem(ctx context.Context, arg dbgen.UpdateTimeSlotTemplateItemParams) (dbgen.TimeSlotTemplateItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.TimeSlotTemplateItem), args.Error(1)
}

func (m *MockQuerier) GetTimeSlotTemplateItemsByTemplateIDExcluding(ctx context.Context, arg dbgen.GetTimeSlotTemplateItemsByTemplateIDExcludingParams) ([]dbgen.TimeSlotTemplateItem, error) {
	args := m.Called(ctx, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dbgen.TimeSlotTemplateItem), args.Error(1)
}

func (m *MockQuerier) DeleteTimeSlotTemplate(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) DeleteTimeSlotTemplateItem(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Store related mock methods
func (m *MockQuerier) CheckStoreNameExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CreateStore(ctx context.Context, arg dbgen.CreateStoreParams) (dbgen.Store, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.Store), args.Error(1)
}
