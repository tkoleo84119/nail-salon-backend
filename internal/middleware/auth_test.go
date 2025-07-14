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
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
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

func TestRequireRoles_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		staffContext := common.StaffContext{
			UserID:   123,
			Username: "testuser",
			Role:     staff.RoleAdmin,
			StoreList: []common.Store{
				{ID: 1, Name: "Store 1"},
			},
		}
		c.Set(UserContextKey, staffContext)
		c.Next()
	})
	router.Use(RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestRequireRoles_InsufficientPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		staffContext := common.StaffContext{
			UserID:   123,
			Username: "testuser",
			Role:     staff.RoleStylist,
			StoreList: []common.Store{
				{ID: 1, Name: "Store 1"},
			},
		}
		c.Set(UserContextKey, staffContext)
		c.Next()
	})
	router.Use(RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "權限不足")
}

func TestRequireRoles_NoStaffContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequireRoles(staff.RoleSuperAdmin, staff.RoleAdmin))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "未找到使用者認證資訊")
}

func TestRequireSuperAdmin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		staffContext := common.StaffContext{
			UserID:   123,
			Username: "testuser",
			Role:     staff.RoleSuperAdmin,
			StoreList: []common.Store{
				{ID: 1, Name: "Store 1"},
			},
		}
		c.Set(UserContextKey, staffContext)
		c.Next()
	})
	router.Use(RequireSuperAdmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireSuperAdmin_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		staffContext := common.StaffContext{
			UserID:   123,
			Username: "testuser",
			Role:     staff.RoleAdmin,
			StoreList: []common.Store{
				{ID: 1, Name: "Store 1"},
			},
		}
		c.Set(UserContextKey, staffContext)
		c.Next()
	})
	router.Use(RequireSuperAdmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireAdminRoles_Success(t *testing.T) {
	testCases := []struct {
		name string
		role string
	}{
		{"SuperAdmin", staff.RoleSuperAdmin},
		{"Admin", staff.RoleAdmin},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				staffContext := common.StaffContext{
					UserID:   123,
					Username: "testuser",
					Role:     tc.role,
					StoreList: []common.Store{
						{ID: 1, Name: "Store 1"},
					},
				}
				c.Set(UserContextKey, staffContext)
				c.Next()
			})
			router.Use(RequireAdminRoles())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRequireManagerOrAbove_Success(t *testing.T) {
	testCases := []struct {
		name string
		role string
	}{
		{"SuperAdmin", staff.RoleSuperAdmin},
		{"Admin", staff.RoleAdmin},
		{"Manager", staff.RoleManager},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				staffContext := common.StaffContext{
					UserID:   123,
					Username: "testuser",
					Role:     tc.role,
					StoreList: []common.Store{
						{ID: 1, Name: "Store 1"},
					},
				}
				c.Set(UserContextKey, staffContext)
				c.Next()
			})
			router.Use(RequireManagerOrAbove())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestRequireAnyStaffRole_Success(t *testing.T) {
	testCases := []struct {
		name string
		role string
	}{
		{"SuperAdmin", staff.RoleSuperAdmin},
		{"Admin", staff.RoleAdmin},
		{"Manager", staff.RoleManager},
		{"Stylist", staff.RoleStylist},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				staffContext := common.StaffContext{
					UserID:   123,
					Username: "testuser",
					Role:     tc.role,
					StoreList: []common.Store{
						{ID: 1, Name: "Store 1"},
					},
				}
				c.Set(UserContextKey, staffContext)
				c.Next()
			})
			router.Use(RequireAnyStaffRole())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestHasRequiredRole(t *testing.T) {
	testCases := []struct {
		name         string
		userRole     string
		allowedRoles []string
		expected     bool
	}{
		{
			name:         "Role exists in allowed list",
			userRole:     staff.RoleAdmin,
			allowedRoles: []string{staff.RoleSuperAdmin, staff.RoleAdmin, staff.RoleManager},
			expected:     true,
		},
		{
			name:         "Role does not exist in allowed list",
			userRole:     staff.RoleStylist,
			allowedRoles: []string{staff.RoleSuperAdmin, staff.RoleAdmin},
			expected:     false,
		},
		{
			name:         "Empty allowed roles",
			userRole:     staff.RoleAdmin,
			allowedRoles: []string{},
			expected:     false,
		},
		{
			name:         "Single role match",
			userRole:     staff.RoleSuperAdmin,
			allowedRoles: []string{staff.RoleSuperAdmin},
			expected:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := hasRequiredRole(tc.userRole, tc.allowedRoles)
			assert.Equal(t, tc.expected, result)
		})
	}
}
