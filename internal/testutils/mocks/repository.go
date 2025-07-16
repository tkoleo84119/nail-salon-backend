package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
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