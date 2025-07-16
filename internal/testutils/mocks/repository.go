package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
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

func (m *MockStylistRepository) UpdateStylist(ctx context.Context, staffUserID int64, req stylist.UpdateStylistRequest) (*stylist.UpdateStylistResponse, error) {
	args := m.Called(ctx, staffUserID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stylist.UpdateStylistResponse), args.Error(1)
}