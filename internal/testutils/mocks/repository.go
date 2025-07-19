package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
)

// MockStaffUserRepository mocks the sqlx repository for testing
type MockStaffUserRepository struct {
	mock.Mock
}

// Ensure MockStaffUserRepository implements the interface
var _ sqlxRepo.StaffUserRepositoryInterface = (*MockStaffUserRepository)(nil)

// NewMockStaffUserRepository creates a new instance of MockStaffUserRepository
func NewMockStaffUserRepository() *MockStaffUserRepository {
	return &MockStaffUserRepository{}
}

func (m *MockStaffUserRepository) UpdateStaffUser(ctx context.Context, id int64, req staff.UpdateStaffRequest) (*staff.UpdateStaffResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.UpdateStaffResponse), args.Error(1)
}

func (m *MockStaffUserRepository) UpdateMyStaff(ctx context.Context, id int64, req staff.UpdateMyStaffRequest) (*staff.UpdateMyStaffResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*staff.UpdateMyStaffResponse), args.Error(1)
}

// MockStylistRepository mocks the sqlx repository for testing
type MockStylistRepository struct {
	mock.Mock
}

// Ensure MockStylistRepository implements the interface
var _ sqlxRepo.StylistRepositoryInterface = (*MockStylistRepository)(nil)

// NewMockStylistRepository creates a new instance of MockStylistRepository
func NewMockStylistRepository() *MockStylistRepository {
	return &MockStylistRepository{}
}

func (m *MockStylistRepository) UpdateStylist(ctx context.Context, staffUserID int64, req stylist.UpdateMyStylistRequest) (*stylist.UpdateMyStylistResponse, error) {
	args := m.Called(ctx, staffUserID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stylist.UpdateMyStylistResponse), args.Error(1)
}

// MockTimeSlotRepository is a mock implementation of TimeSlotRepositoryInterface
type MockTimeSlotRepository struct {
	mock.Mock
}

var _ sqlxRepo.TimeSlotRepositoryInterface = (*MockTimeSlotRepository)(nil)

func NewMockTimeSlotRepository() *MockTimeSlotRepository {
	return &MockTimeSlotRepository{}
}

func (m *MockTimeSlotRepository) UpdateTimeSlot(ctx context.Context, timeSlotID int64, req schedule.UpdateTimeSlotRequest) (*schedule.UpdateTimeSlotResponse, error) {
	args := m.Called(ctx, timeSlotID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.UpdateTimeSlotResponse), args.Error(1)
}

// MockUpdateTimeSlotTemplateRepository implements the UpdateTimeSlotTemplateRepositoryInterface for testing
type MockUpdateTimeSlotTemplateRepository struct {
	mock.Mock
}

var _ sqlxRepo.TimeSlotTemplateRepositoryInterface = (*MockUpdateTimeSlotTemplateRepository)(nil)

func NewMockUpdateTimeSlotTemplateRepository() *MockUpdateTimeSlotTemplateRepository {
	return &MockUpdateTimeSlotTemplateRepository{}
}

func (m *MockUpdateTimeSlotTemplateRepository) UpdateTimeSlotTemplate(ctx context.Context, templateID int64, req timeSlotTemplate.UpdateTimeSlotTemplateRequest) (*timeSlotTemplate.UpdateTimeSlotTemplateResponse, error) {
	args := m.Called(ctx, templateID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*timeSlotTemplate.UpdateTimeSlotTemplateResponse), args.Error(1)
}

// MockStoreRepository mocks the sqlx repository for testing
type MockStoreRepository struct {
	mock.Mock
}

// Ensure MockStoreRepository implements the interface
var _ sqlxRepo.StoreRepositoryInterface = (*MockStoreRepository)(nil)

// NewMockStoreRepository creates a new instance of MockStoreRepository
func NewMockStoreRepository() *MockStoreRepository {
	return &MockStoreRepository{}
}

func (m *MockStoreRepository) UpdateStore(ctx context.Context, storeID int64, req store.UpdateStoreRequest) (*store.UpdateStoreResponse, error) {
	args := m.Called(ctx, storeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.UpdateStoreResponse), args.Error(1)
}
