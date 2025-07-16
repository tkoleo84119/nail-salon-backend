package auth

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/mocks"
	"github.com/tkoleo84119/nail-salon-backend/internal/testutils/setup"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)


func TestLoginService_Login_Success(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLoginService(mockQuerier, env.Config.JWT)

	// Hash password for test user
	hashedPassword, err := utils.HashPassword("testpassword")
	require.NoError(t, err)

	// Set up mock expectations
	testUser := dbgen.StaffUser{
		ID:           123,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         staff.RoleAdmin,
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	}

	storeAccess := []dbgen.GetStaffUserStoreAccessRow{
		{StoreID: 1, StoreName: "Store 1"},
		{StoreID: 2, StoreName: "Store 2"},
	}

	tokenRow := dbgen.CreateStaffUserTokenRow{
		ID:        456,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	mockQuerier.On("GetStaffUserByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockQuerier.On("GetStaffUserStoreAccess", mock.Anything, int64(123)).Return(storeAccess, nil)
	mockQuerier.On("CreateStaffUserToken", mock.Anything, mock.AnythingOfType("dbgen.CreateStaffUserTokenParams")).Return(tokenRow, nil)

	// Create request and context
	req := auth.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	loginCtx := auth.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call service
	response, err := service.Login(context.Background(), req, loginCtx)

	// Assert response
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, 3600, response.ExpiresIn)
	assert.Equal(t, "123", response.User.ID)
	assert.Equal(t, "testuser", response.User.Username)
	assert.Equal(t, staff.RoleAdmin, response.User.Role)
	assert.Equal(t, []common.Store{
		{ID: "1", Name: "Store 1"},
		{ID: "2", Name: "Store 2"},
	}, response.User.StoreList)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestLoginService_Login_SuperAdmin(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLoginService(mockQuerier, env.Config.JWT)

	// Hash password for test user
	hashedPassword, err := utils.HashPassword("adminpassword")
	require.NoError(t, err)

	// Set up mock expectations for SUPER_ADMIN
	testUser := dbgen.StaffUser{
		ID:           999,
		Username:     "superadmin",
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		Role:         staff.RoleSuperAdmin,
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	}

	allStores := []dbgen.GetAllActiveStoresRow{
		{ID: 1, Name: "Store 1"},
		{ID: 2, Name: "Store 2"},
		{ID: 3, Name: "Store 3"},
	}

	tokenRow := dbgen.CreateStaffUserTokenRow{
		ID:        789,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	mockQuerier.On("GetStaffUserByUsername", mock.Anything, "superadmin").Return(testUser, nil)
	mockQuerier.On("GetAllActiveStores", mock.Anything).Return(allStores, nil)
	mockQuerier.On("CreateStaffUserToken", mock.Anything, mock.AnythingOfType("dbgen.CreateStaffUserTokenParams")).Return(tokenRow, nil)

	// Create request and context
	req := auth.LoginRequest{
		Username: "superadmin",
		Password: "adminpassword",
	}

	loginCtx := auth.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call service
	response, err := service.Login(context.Background(), req, loginCtx)

	// Assert response
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, staff.RoleSuperAdmin, response.User.Role)
	assert.Equal(t, []common.Store{
		{ID: "1", Name: "Store 1"},
		{ID: "2", Name: "Store 2"},
		{ID: "3", Name: "Store 3"},
	}, response.User.StoreList)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestLoginService_Login_InvalidCredentials(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLoginService(mockQuerier, env.Config.JWT)

	// Set up mock expectations - user not found
	mockQuerier.On("GetStaffUserByUsername", mock.Anything, "nonexistent").Return(dbgen.StaffUser{}, assert.AnError)

	// Create request
	req := auth.LoginRequest{
		Username: "nonexistent",
		Password: "password",
	}

	loginCtx := auth.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call service
	response, err := service.Login(context.Background(), req, loginCtx)

	// Assert error
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_INVALID_CREDENTIALS")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestLoginService_Login_WrongPassword(t *testing.T) {
	env := setup.SetupTestEnvironmentForService(t)
	defer env.Cleanup()

	// Create mock querier
	mockQuerier := mocks.NewMockQuerier()
	service := NewLoginService(mockQuerier, env.Config.JWT)

	// Hash a different password
	hashedPassword, err := utils.HashPassword("correctpassword")
	require.NoError(t, err)

	// Set up mock expectations
	testUser := dbgen.StaffUser{
		ID:           123,
		Username:     "testuser",
		PasswordHash: hashedPassword,
		Role:         staff.RoleAdmin,
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	}

	mockQuerier.On("GetStaffUserByUsername", mock.Anything, "testuser").Return(testUser, nil)

	// Create request with wrong password
	req := auth.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	loginCtx := auth.LoginContext{
		UserAgent: "test-agent",
		IPAddress: "127.0.0.1",
		Timestamp: time.Now(),
	}

	// Call service
	response, err := service.Login(context.Background(), req, loginCtx)

	// Assert error
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AUTH_INVALID_CREDENTIALS")

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}
