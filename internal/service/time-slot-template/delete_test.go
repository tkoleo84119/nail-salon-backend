package timeSlotTemplate

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
)

func TestDeleteTimeSlotTemplateService_DeleteTimeSlotTemplate_Success(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateService(mockQuerier)

	ctx := context.Background()
	templateID := "6000000011"
	templateIDInt := int64(6000000011)
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Mock existing template
	existingTemplate := dbgen.TimeSlotTemplate{
		ID:   templateIDInt,
		Name: "Test Template",
		Note: pgtype.Text{String: "Test note", Valid: true},
	}

	mockQuerier.On("GetTimeSlotTemplateByID", ctx, templateIDInt).Return(existingTemplate, nil)
	mockQuerier.On("DeleteTimeSlotTemplate", ctx, templateIDInt).Return(nil)

	response, err := service.DeleteTimeSlotTemplate(ctx, templateID, staffContext)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, []string{templateID}, response.Deleted)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateService_DeleteTimeSlotTemplate_InvalidTemplateID(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateService(mockQuerier)

	ctx := context.Background()
	templateID := "invalid"
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	response, err := service.DeleteTimeSlotTemplate(ctx, templateID, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.ValInputValidationFailed, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateService_DeleteTimeSlotTemplate_TemplateNotFound(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateService(mockQuerier)

	ctx := context.Background()
	templateID := "6000000011"
	templateIDInt := int64(6000000011)
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	mockQuerier.On("GetTimeSlotTemplateByID", ctx, templateIDInt).Return(dbgen.TimeSlotTemplate{}, errors.New("not found"))

	response, err := service.DeleteTimeSlotTemplate(ctx, templateID, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.TimeSlotTemplateNotFound, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}

func TestDeleteTimeSlotTemplateService_DeleteTimeSlotTemplate_DatabaseError(t *testing.T) {
	mockQuerier := mocks.NewMockQuerier()
	service := NewDeleteTimeSlotTemplateService(mockQuerier)

	ctx := context.Background()
	templateID := "6000000011"
	templateIDInt := int64(6000000011)
	staffContext := common.StaffContext{
		UserID: "11111",
		Role:   staff.RoleAdmin,
	}

	// Mock existing template
	existingTemplate := dbgen.TimeSlotTemplate{
		ID:   templateIDInt,
		Name: "Test Template",
		Note: pgtype.Text{String: "Test note", Valid: true},
	}

	mockQuerier.On("GetTimeSlotTemplateByID", ctx, templateIDInt).Return(existingTemplate, nil)
	mockQuerier.On("DeleteTimeSlotTemplate", ctx, templateIDInt).Return(errors.New("database error"))

	response, err := service.DeleteTimeSlotTemplate(ctx, templateID, staffContext)

	assert.Nil(t, response)
	assert.Error(t, err)

	serviceErr, ok := err.(*errorCodes.ServiceError)
	assert.True(t, ok)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockQuerier.AssertExpectations(t)
}