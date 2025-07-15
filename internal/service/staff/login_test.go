package staff

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type MockQuerier struct {
	mock.Mock
}

var _ dbgen.Querier = (*MockQuerier)(nil)

func (m *MockQuerier) GetStaffUserByUsername(ctx context.Context, username string) (dbgen.StaffUser, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) GetStaffUserStoreAccess(ctx context.Context, staffUserID int64) ([]dbgen.GetStaffUserStoreAccessRow, error) {
	args := m.Called(ctx, staffUserID)
	return args.Get(0).([]dbgen.GetStaffUserStoreAccessRow), args.Error(1)
}

func (m *MockQuerier) GetAllActiveStores(ctx context.Context) ([]dbgen.GetAllActiveStoresRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dbgen.GetAllActiveStoresRow), args.Error(1)
}

func (m *MockQuerier) GetStaffUserByID(ctx context.Context, userID int64) (dbgen.StaffUser, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) CreateStaffUserToken(ctx context.Context, arg dbgen.CreateStaffUserTokenParams) (dbgen.CreateStaffUserTokenRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.CreateStaffUserTokenRow), args.Error(1)
}

func (m *MockQuerier) CheckStaffUserExists(ctx context.Context, arg dbgen.CheckStaffUserExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CheckStoresExistAndActive(ctx context.Context, storeIDs []int64) (dbgen.CheckStoresExistAndActiveRow, error) {
	args := m.Called(ctx, storeIDs)
	return args.Get(0).(dbgen.CheckStoresExistAndActiveRow), args.Error(1)
}

func (m *MockQuerier) CreateStaffUser(ctx context.Context, arg dbgen.CreateStaffUserParams) (dbgen.CreateStaffUserRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.CreateStaffUserRow), args.Error(1)
}

func (m *MockQuerier) CreateStaffUserStoreAccess(ctx context.Context, arg dbgen.CreateStaffUserStoreAccessParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetStoresByIDs(ctx context.Context, storeIDs []int64) ([]dbgen.GetStoresByIDsRow, error) {
	args := m.Called(ctx, storeIDs)
	return args.Get(0).([]dbgen.GetStoresByIDsRow), args.Error(1)
}

func (m *MockQuerier) BatchCreateStaffUserStoreAccess(ctx context.Context, arg dbgen.BatchCreateStaffUserStoreAccessParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) GetStoreByID(ctx context.Context, id int64) (dbgen.GetStoreByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dbgen.GetStoreByIDRow), args.Error(1)
}

func (m *MockQuerier) CheckStoreAccessExists(ctx context.Context, arg dbgen.CheckStoreAccessExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func setupTestEnvironment(t *testing.T) func() {
	// Set up JWT secret for testing
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-staff-service-testing")

	// Initialize snowflake for testing
	err := utils.InitSnowflake(1)
	require.NoError(t, err)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	err = errorManager.LoadFromFile("../../errors/errors.yaml")
	require.NoError(t, err)

	return func() {
		os.Setenv("JWT_SECRET", originalSecret)
	}
}

func TestLoginService_Login_Success(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockQuerier)
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key-for-staff-service-testing",
		ExpiryHours: 1,
	}
	service := NewLoginService(mockQuerier, jwtConfig)

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
	req := staff.LoginRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	loginCtx := staff.LoginContext{
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
		{ID: 1, Name: "Store 1"},
		{ID: 2, Name: "Store 2"},
	}, response.User.StoreList)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestLoginService_Login_SuperAdmin(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockQuerier)
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key-for-staff-service-testing",
		ExpiryHours: 1,
	}
	service := NewLoginService(mockQuerier, jwtConfig)

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
	req := staff.LoginRequest{
		Username: "superadmin",
		Password: "adminpassword",
	}

	loginCtx := staff.LoginContext{
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
		{ID: 1, Name: "Store 1"},
		{ID: 2, Name: "Store 2"},
		{ID: 3, Name: "Store 3"},
	}, response.User.StoreList)

	// Verify all expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestLoginService_Login_InvalidCredentials(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockQuerier)
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key-for-staff-service-testing",
		ExpiryHours: 1,
	}
	service := NewLoginService(mockQuerier, jwtConfig)

	// Set up mock expectations - user not found
	mockQuerier.On("GetStaffUserByUsername", mock.Anything, "nonexistent").Return(dbgen.StaffUser{}, assert.AnError)

	// Create request
	req := staff.LoginRequest{
		Username: "nonexistent",
		Password: "password",
	}

	loginCtx := staff.LoginContext{
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
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create mock querier
	mockQuerier := new(MockQuerier)
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key-for-staff-service-testing",
		ExpiryHours: 1,
	}
	service := NewLoginService(mockQuerier, jwtConfig)

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
	req := staff.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	loginCtx := staff.LoginContext{
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
