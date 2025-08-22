package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ErrorDefinition represents a single error definition
type ErrorDefinition struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// ErrorManager manages error definitions loaded from YAML
type ErrorManager struct {
	errors    map[string]ErrorDefinition
	codeToKey map[string]string
	mu        sync.RWMutex
}

var (
	manager *ErrorManager
	once    sync.Once
)

// GetCode retrieves error code by constant name
func (em *ErrorManager) GetCode(constantName string) string {
	if def, exists := em.GetError(constantName); exists {
		return def.Code
	}
	return "E9999" // default unknown error code
}

// GetManager returns the singleton error manager instance
func GetManager() *ErrorManager {
	once.Do(func() {
		manager = &ErrorManager{
			errors: make(map[string]ErrorDefinition),
		}
	})
	return manager
}

// LoadFromFile loads error definitions from JSON file
func (em *ErrorManager) LoadFromFile(filepath string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read error file: %w", err)
	}

	var categories map[string]map[string]ErrorDefinition
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// initialize reverse mapping
	em.codeToKey = make(map[string]string)

	// Flatten the categories into a single map
	for _, errorDefs := range categories {
		for key, def := range errorDefs {
			em.errors[key] = def
			em.codeToKey[def.Code] = key
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

// GetErrorResponse creates a standardized error response
func (em *ErrorManager) GetErrorResponse(errors []ErrorItem) map[string]interface{} {
	return map[string]interface{}{
		"errors": errors,
	}
}

// CreateErrorItem creates an error item from a constant name and optional field and parameters
func (em *ErrorManager) CreateErrorItem(constantName string, field string, params map[string]string) ErrorItem {
	def, exists := em.GetError(constantName)
	if !exists {
		return ErrorItem{
			Code:    "E9999",
			Message: "未知錯誤",
		}
	}

	message := def.Message
	for key, value := range params {
		placeholder := "{" + key + "}"
		message = strings.ReplaceAll(message, placeholder, value)
	}

	item := ErrorItem{
		Code:    def.Code,
		Message: message,
	}

	if field != "" {
		item.Field = field
	}

	return item
}

// GetStatusByCode retrieves HTTP status code by error code
func (em *ErrorManager) GetStatusByCode(errorCode string) int {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if key, exists := em.codeToKey[errorCode]; exists {
		if def, exists := em.errors[key]; exists {
			return def.Status
		}
	}
	return 500
}

// isReleaseMode checks if the application is running in release mode
func (em *ErrorManager) isDebugMode() bool {
	return gin.Mode() == gin.DebugMode
}
