package timeSlotTemplate

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	name := "Updated Template"
	note := "Updated note"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
		Note: &note,
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Original Template",
	}, nil)

	// Mock repository update
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateResponse{
		ID:   "6000000011",
		Name: "Updated Template",
		Note: "Updated note",
	}
	mockRepo.On("UpdateTimeSlotTemplate", ctx, int64(6000000011), request).Return(expectedResponse, nil)

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "6000000011", response.ID)
	assert.Equal(t, "Updated Template", response.Name)
	assert.Equal(t, "Updated note", response.Note)

	mockQuerier.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_UpdateNameOnly(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleManager,
	}

	templateID := "6000000011"
	name := "Updated Template"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Original Template",
	}, nil)

	// Mock repository update
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateResponse{
		ID:   "6000000011",
		Name: "Updated Template",
		Note: "Original note",
	}
	mockRepo.On("UpdateTimeSlotTemplate", ctx, int64(6000000011), request).Return(expectedResponse, nil)

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "6000000011", response.ID)
	assert.Equal(t, "Updated Template", response.Name)
	assert.Equal(t, "Original note", response.Note)

	mockQuerier.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_UpdateNoteOnly(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleSuperAdmin,
	}

	templateID := "6000000011"
	note := "Updated note"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Note: &note,
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Original Template",
	}, nil)

	// Mock repository update
	expectedResponse := &timeSlotTemplate.UpdateTimeSlotTemplateResponse{
		ID:   "6000000011",
		Name: "Original Template",
		Note: "Updated note",
	}
	mockRepo.On("UpdateTimeSlotTemplate", ctx, int64(6000000011), request).Return(expectedResponse, nil)

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "6000000011", response.ID)
	assert.Equal(t, "Original Template", response.Name)
	assert.Equal(t, "Updated note", response.Note)

	mockQuerier.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_AllFieldsEmpty(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		// Both fields are nil
	}

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValAllFieldsEmpty, serviceErr.Code)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_InvalidTemplateID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "invalid"
	name := "Updated Template"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_TemplateNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	name := "Updated Template"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}

	// Mock template not found
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{}, errors.New("not found"))

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestUpdateTimeSlotTemplateService_UpdateTimeSlotTemplate_RepositoryError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	mockRepo := mocks.NewMockUpdateTimeSlotTemplateRepository()
	service := NewUpdateTimeSlotTemplateService(mockQuerier, mockRepo)

	ctx := context.Background()
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	templateID := "6000000011"
	name := "Updated Template"
	request := timeSlotTemplate.UpdateTimeSlotTemplateRequest{
		Name: &name,
	}

	// Mock template exists
	mockQuerier.On("GetTimeSlotTemplateByID", ctx, int64(6000000011)).Return(dbgen.TimeSlotTemplate{
		ID:   6000000011,
		Name: "Original Template",
	}, nil)

	// Mock repository error
	mockRepo.On("UpdateTimeSlotTemplate", ctx, int64(6000000011), request).Return(nil, errors.New("database error"))

	// Call service
	response, err := service.UpdateTimeSlotTemplate(ctx, templateID, request, staffContext)

	// Assertions
	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
