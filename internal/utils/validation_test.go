package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Username string `validate:"required,min=1,max=100"`
	Password string `validate:"required,min=1,max=100"`
}

func TestExtractValidationErrors(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name          string
		input         TestStruct
		expectedKeys  []string
		expectedMsgs  map[string]string
	}{
		{
			name:         "required fields missing",
			input:        TestStruct{},
			expectedKeys: []string{"username", "password"},
			expectedMsgs: map[string]string{
				"username": "帳號為必填項目",
				"password": "密碼為必填項目",
			},
		},
		{
			name: "min length validation",
			input: TestStruct{
				Username: "",
				Password: "",
			},
			expectedKeys: []string{"username", "password"},
			expectedMsgs: map[string]string{
				"username": "帳號為必填項目",
				"password": "密碼為必填項目",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			if err != nil {
				errors := ExtractValidationErrors(err)
				
				// Check that all expected keys are present
				for _, key := range tt.expectedKeys {
					assert.Contains(t, errors, key, "Expected key %s not found", key)
				}
				
				// Check specific error messages
				for key, expectedMsg := range tt.expectedMsgs {
					assert.Equal(t, expectedMsg, errors[key], "Error message mismatch for key %s", key)
				}
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "validation error",
			input:    TestStruct{},
			expected: true,
		},
		{
			name:     "valid struct",
			input:    TestStruct{Username: "test", Password: "test"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			
			if tt.expected {
				assert.Error(t, err)
				assert.True(t, IsValidationError(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}