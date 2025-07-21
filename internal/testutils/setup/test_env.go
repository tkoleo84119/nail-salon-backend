package setup

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TestEnvironment holds test configuration and cleanup functions
type TestEnvironment struct {
	Config       config.Config
	ErrorManager *errorCodes.ErrorManager
	cleanup      func()
}

// Cleanup performs environment cleanup
func (env *TestEnvironment) Cleanup() {
	if env.cleanup != nil {
		env.cleanup()
	}
}

// SetupTestEnvironmentForService sets up test environment for service layer tests
func SetupTestEnvironmentForService(t *testing.T) *TestEnvironment {
	// Store original JWT secret
	originalSecret := os.Getenv("JWT_SECRET")
	
	// Set test JWT secret
	testSecret := "test-secret-key-for-staff-service-testing"
	os.Setenv("JWT_SECRET", testSecret)

	// Initialize snowflake for testing
	err := utils.InitSnowflake(1)
	require.NoError(t, err)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	err = errorManager.LoadFromFile("../../errors/errors.yaml")
	require.NoError(t, err)

	// Create test config
	testConfig := config.Config{
		JWT: config.JWTConfig{
			Secret:      testSecret,
			ExpiryHours: 1,
		},
	}

	// Create cleanup function
	cleanup := func() {
		os.Setenv("JWT_SECRET", originalSecret)
	}

	return &TestEnvironment{
		Config:       testConfig,
		ErrorManager: errorManager,
		cleanup:      cleanup,
	}
}

// SetupTestEnvironmentForHandler sets up test environment for handler layer tests
func SetupTestEnvironmentForHandler(t *testing.T) *TestEnvironment {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("taiwanlandline", utils.ValidateTaiwanLandline)
		v.RegisterValidation("taiwanmobile", utils.ValidateTaiwanMobile)
	}

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	err := errorManager.LoadFromFile("../../errors/errors.yaml")
	require.NoError(t, err)

	// Create test config (minimal for handlers)
	testConfig := config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	return &TestEnvironment{
		Config:       testConfig,
		ErrorManager: errorManager,
		cleanup:      nil, // No cleanup needed for handler tests
	}
}

// SetupTestEnvironmentForMiddleware sets up test environment for middleware tests
func SetupTestEnvironmentForMiddleware(t *testing.T) *TestEnvironment {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load error definitions for testing
	errorManager := errorCodes.GetManager()
	err := errorManager.LoadFromFile("../errors/errors.yaml")
	require.NoError(t, err)

	// Create test config
	testConfig := config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	return &TestEnvironment{
		Config:       testConfig,
		ErrorManager: errorManager,
		cleanup:      nil, // No cleanup needed for middleware tests
	}
}