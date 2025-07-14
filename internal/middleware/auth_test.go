package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
)

type MockQuerier struct {
	mock.Mock
}

var _ dbgen.Querier = (*MockQuerier)(nil)

func (m *MockQuerier) GetStaffUserByID(ctx context.Context, userID int64) (dbgen.StaffUser, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) GetStaffUserByUsername(ctx context.Context, username string) (dbgen.StaffUser, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(dbgen.StaffUser), args.Error(1)
}

func (m *MockQuerier) CreateStaffUserToken(ctx context.Context, arg dbgen.CreateStaffUserTokenParams) (dbgen.CreateStaffUserTokenRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(dbgen.CreateStaffUserTokenRow), args.Error(1)
}

func (m *MockQuerier) GetAllActiveStores(ctx context.Context) ([]dbgen.GetAllActiveStoresRow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dbgen.GetAllActiveStoresRow), args.Error(1)
}

func (m *MockQuerier) GetStaffUserStoreAccess(ctx context.Context, staffUserID int64) ([]dbgen.GetStaffUserStoreAccessRow, error) {
	args := m.Called(ctx, staffUserID)
	return args.Get(0).([]dbgen.GetStaffUserStoreAccessRow), args.Error(1)
}

func TestJWTAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(MockQuerier)
	cfg := config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	router := gin.New()
	router.Use(JWTAuth(cfg, mockDB))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "access_token 缺失")
}

func TestJWTAuth_InvalidTokenFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDB := new(MockQuerier)
	cfg := config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	router := gin.New()
	router.Use(JWTAuth(cfg, mockDB))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "access_token 格式錯誤")
}

func TestGetStaffFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	staffContext := common.StaffContext{
		UserID:   123,
		Username: "testuser",
		Role:     "ADMIN",
		StoreList: []common.Store{
			{ID: 1, Name: "Store 1"},
		},
	}

	c.Set(UserContextKey, staffContext)

	result, exists := GetStaffFromContext(c)

	assert.True(t, exists)
	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.UserID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "ADMIN", result.Role)
}

func TestGetStaffFromContext_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	result, exists := GetStaffFromContext(c)

	assert.False(t, exists)
	assert.Nil(t, result)
}
