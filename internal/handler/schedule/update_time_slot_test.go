package schedule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestNewUpdateTimeSlotHandler(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockTimeSlotRepo := &MockTimeSlotRepository{}
	service := scheduleService.NewUpdateTimeSlotService(mockQuerier, mockTimeSlotRepo)
	handler := NewUpdateTimeSlotHandler(service)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}

// MockTimeSlotRepository is a mock implementation of TimeSlotRepositoryInterface
type MockTimeSlotRepository struct {
	mock.Mock
}

func (m *MockTimeSlotRepository) UpdateTimeSlot(ctx context.Context, timeSlotID int64, req schedule.UpdateTimeSlotRequest) (*schedule.UpdateTimeSlotResponse, error) {
	args := m.Called(ctx, timeSlotID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schedule.UpdateTimeSlotResponse), args.Error(1)
}