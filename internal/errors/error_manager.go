package errors

import (
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// ErrorDefinition represents a single error definition
type ErrorDefinition struct {
	Message string `yaml:"message"`
	Status  int    `yaml:"status"`
}

// ErrorManager manages error definitions loaded from YAML
type ErrorManager struct {
	errors map[string]ErrorDefinition
	mu     sync.RWMutex
}

var (
	manager *ErrorManager
	once    sync.Once
)

// GetManager returns the singleton error manager instance
func GetManager() *ErrorManager {
	once.Do(func() {
		manager = &ErrorManager{
			errors: make(map[string]ErrorDefinition),
		}
	})
	return manager
}

// LoadFromFile loads error definitions from YAML file
func (em *ErrorManager) LoadFromFile(filepath string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read error file: %w", err)
	}

	var categories map[string]map[string]ErrorDefinition
	if err := yaml.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Flatten the categories into a single map
	for _, errorDefs := range categories {
		for code, def := range errorDefs {
			em.errors[code] = def
		}
	}

	return nil
}

// GetError retrieves error definition by code
func (em *ErrorManager) GetError(code string) (ErrorDefinition, bool) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	def, exists := em.errors[code]
	return def, exists
}

// GetMessage retrieves error message by code
func (em *ErrorManager) GetMessage(code string) string {
	if def, exists := em.GetError(code); exists {
		return def.Message
	}
	return "未知錯誤"
}

// GetStatus retrieves HTTP status code by error code
func (em *ErrorManager) GetStatus(code string) int {
	if def, exists := em.GetError(code); exists {
		return def.Status
	}
	return 500
}

// GetErrorResponse creates a standardized error response with optional development details
func (em *ErrorManager) GetErrorResponse(code string, fieldErrors map[string]string, details ...string) map[string]interface{} {
	message := em.GetMessage(code)

	response := map[string]interface{}{
		"message": message,
	}

	if len(fieldErrors) > 0 {
		response["errors"] = fieldErrors
	}

	// Add development details if provided and in debug mode
	if len(details) > 0 && details[0] != "" && em.isDebugMode() {
		response["dev_details"] = details[0]
	}

	return response
}

// isReleaseMode checks if the application is running in release mode
func (em *ErrorManager) isDebugMode() bool {
	return gin.Mode() == gin.DebugMode
}
