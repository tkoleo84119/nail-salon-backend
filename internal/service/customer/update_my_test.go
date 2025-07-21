package customer

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
)

func TestUpdateMyCustomerService_UpdateMyCustomer_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)
	name := "王小美"
	phone := "0912345678"
	birthday := "1992-02-29"
	city := "台北市"
	favoriteShapes := []string{"方形"}
	favoriteColors := []string{"粉色"}
	favoriteStyles := []string{"法式"}
	isIntrovert := true
	customerNote := "容易指緣乾裂"

	req := customer.UpdateMyCustomerRequest{
		Name:           &name,
		Phone:          &phone,
		Birthday:       &birthday,
		City:           &city,
		FavoriteShapes: &favoriteShapes,
		FavoriteColors: &favoriteColors,
		FavoriteStyles: &favoriteStyles,
		IsIntrovert:    &isIntrovert,
		CustomerNote:   &customerNote,
	}

	existingCustomer := dbgen.Customer{
		ID:    customerID,
		Name:  "Old Name",
		Phone: "0987654321",
	}

	expectedResponse := &customer.UpdateMyCustomerResponse{
		ID:             "1000000001",
		Name:           name,
		Phone:          phone,
		Birthday:       &birthday,
		City:           &city,
		FavoriteShapes: &favoriteShapes,
		FavoriteColors: &favoriteColors,
		FavoriteStyles: &favoriteStyles,
		IsIntrovert:    &isIntrovert,
		CustomerNote:   &customerNote,
	}

	mockDB.On("GetCustomerByID", ctx, customerID).Return(existingCustomer, nil)
	mockCustomerRepo.On("UpdateMyCustomer", ctx, customerID, req).Return(expectedResponse, nil)

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.ID, result.ID)
	assert.Equal(t, expectedResponse.Name, result.Name)
	assert.Equal(t, expectedResponse.Phone, result.Phone)
	assert.Equal(t, expectedResponse.Birthday, result.Birthday)
	assert.Equal(t, expectedResponse.City, result.City)
	assert.Equal(t, expectedResponse.FavoriteShapes, result.FavoriteShapes)
	assert.Equal(t, expectedResponse.FavoriteColors, result.FavoriteColors)
	assert.Equal(t, expectedResponse.FavoriteStyles, result.FavoriteStyles)
	assert.Equal(t, expectedResponse.IsIntrovert, result.IsIntrovert)
	assert.Equal(t, expectedResponse.CustomerNote, result.CustomerNote)

	mockDB.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
}

func TestUpdateMyCustomerService_UpdateMyCustomer_NoFieldsToUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)

	req := customer.UpdateMyCustomerRequest{} // Empty request

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.Error(t, err)
	assert.Nil(t, result)

	serviceErr := err.(*errorCodes.ServiceError)
	assert.Equal(t, errorCodes.ValAllFieldsEmpty, serviceErr.Code)

	// Should not call repository methods
	mockDB.AssertNotCalled(t, "GetCustomerByID")
	mockCustomerRepo.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerService_UpdateMyCustomer_InvalidBirthdayFormat(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)
	invalidBirthday := "1992/02/29" // Wrong format, should be yyyy-MM-dd

	req := customer.UpdateMyCustomerRequest{
		Birthday: &invalidBirthday,
	}

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.Error(t, err)
	assert.Nil(t, result)

	serviceErr := err.(*errorCodes.ServiceError)
	assert.Equal(t, errorCodes.ValDateFormatInvalid, serviceErr.Code)

	// Should not call repository methods
	mockDB.AssertNotCalled(t, "GetCustomerByID")
	mockCustomerRepo.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerService_UpdateMyCustomer_CustomerNotFound(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)
	name := "王小美"

	req := customer.UpdateMyCustomerRequest{
		Name: &name,
	}

	mockDB.On("GetCustomerByID", ctx, customerID).Return(dbgen.Customer{}, sql.ErrNoRows)

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.Error(t, err)
	assert.Nil(t, result)

	serviceErr := err.(*errorCodes.ServiceError)
	assert.Equal(t, errorCodes.CustomerNotFound, serviceErr.Code)

	mockDB.AssertExpectations(t)
	mockCustomerRepo.AssertNotCalled(t, "UpdateMyCustomer")
}

func TestUpdateMyCustomerService_UpdateMyCustomer_RepositoryError(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)
	name := "王小美"

	req := customer.UpdateMyCustomerRequest{
		Name: &name,
	}

	existingCustomer := dbgen.Customer{
		ID:   customerID,
		Name: "Old Name",
	}

	mockDB.On("GetCustomerByID", ctx, customerID).Return(existingCustomer, nil)
	mockCustomerRepo.On("UpdateMyCustomer", ctx, customerID, req).Return(nil, assert.AnError)

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.Error(t, err)
	assert.Nil(t, result)

	serviceErr := err.(*errorCodes.ServiceError)
	assert.Equal(t, errorCodes.SysDatabaseError, serviceErr.Code)

	mockDB.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
}

func TestUpdateMyCustomerService_UpdateMyCustomer_OnlyNameUpdate(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	mockCustomerRepo := mocks.NewMockCustomerRepository()
	mockDB := mocks.NewMockQuerier()
	service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

	ctx := context.Background()
	customerID := int64(1000000001)
	name := "新名字"

	req := customer.UpdateMyCustomerRequest{
		Name: &name,
	}

	existingCustomer := dbgen.Customer{
		ID:    customerID,
		Name:  "舊名字",
		Phone: "0912345678",
	}

	expectedResponse := &customer.UpdateMyCustomerResponse{
		ID:    "1000000001",
		Name:  name,
		Phone: "0912345678",
	}

	mockDB.On("GetCustomerByID", ctx, customerID).Return(existingCustomer, nil)
	mockCustomerRepo.On("UpdateMyCustomer", ctx, customerID, req).Return(expectedResponse, nil)

	result, err := service.UpdateMyCustomer(ctx, customerID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Name, result.Name)

	mockDB.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
}

func TestUpdateMyCustomerService_UpdateMyCustomer_ValidBirthdayFormats(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	testCases := []struct {
		name     string
		birthday string
	}{
		{"Valid leap year", "2024-02-29"},
		{"Valid regular date", "1990-12-25"},
		{"Valid start of year", "2000-01-01"},
		{"Valid end of year", "2000-12-31"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCustomerRepo := mocks.NewMockCustomerRepository()
			mockDB := mocks.NewMockQuerier()
			service := NewUpdateMyCustomerService(mockDB, mockCustomerRepo)

			ctx := context.Background()
			customerID := int64(1000000001)

			req := customer.UpdateMyCustomerRequest{
				Birthday: &tc.birthday,
			}

			existingCustomer := dbgen.Customer{
				ID: customerID,
			}

			expectedResponse := &customer.UpdateMyCustomerResponse{
				ID:       "1000000001",
				Name:     "Test Name",
				Phone:    "0912345678",
				Birthday: &tc.birthday,
			}

			mockDB.On("GetCustomerByID", ctx, customerID).Return(existingCustomer, nil)
			mockCustomerRepo.On("UpdateMyCustomer", ctx, customerID, req).Return(expectedResponse, nil)

			result, err := service.UpdateMyCustomer(ctx, customerID, req)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, &tc.birthday, result.Birthday)

			mockDB.AssertExpectations(t)
			mockCustomerRepo.AssertExpectations(t)
		})
	}
}
