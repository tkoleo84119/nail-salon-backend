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
		name         string
		input        TestStruct
		expectedKeys []string
		expectedMsgs map[string]string
	}{
		{
			name:         "required fields missing",
			input:        TestStruct{},
			expectedKeys: []string{"username", "password"},
			expectedMsgs: map[string]string{
				"username": "username為必填項目",
				"password": "password為必填項目",
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
				"username": "username為必填項目",
				"password": "password為必填項目",
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

func TestValidateTaiwanLandline(t *testing.T) {
	// Create a validator instance and register our custom validation
	validate := validator.New()
	validate.RegisterValidation("taiwanlandline", ValidateTaiwanLandline)

	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		// Valid cases
		{"Valid Taipei landline", "02-12345678", true},
		{"Valid Taichung landline", "04-23456789", true},
		{"Valid Kaohsiung landline", "07-7654321", true},
		{"Valid Taitung landline", "089-123456", true},
		{"Valid 7-digit number", "03-1234567", true},
		{"Valid 8-digit number", "06-12345678", true},
		{"Empty string should pass", "", true},

		// Invalid cases
		{"Missing area code", "12345678", false},
		{"Wrong area code", "01-12345678", false},
		{"Wrong area code 09", "09-12345678", false},
		{"Missing dash", "0212345678", false},
		{"Too short number", "02-123456", false},
		{"Too long number", "02-123456789", false},
		{"Invalid characters", "02-abcd1234", false},
		{"Wrong format with spaces", "02 12345678", false},
		{"Invalid area code 089 with wrong number length", "089-1234567", false},
		{"Invalid area code 089 with too short number", "089-12345", false},
	}

	type TestStructPhone struct {
		Phone string `validate:"taiwanlandline"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := TestStructPhone{Phone: tt.phone}
			err := validate.Struct(testStruct)

			if tt.expected {
				assert.NoError(t, err, "Expected phone %s to be valid", tt.phone)
			} else {
				assert.Error(t, err, "Expected phone %s to be invalid", tt.phone)
			}
		})
	}
}

func TestExtractValidationErrors_TaiwanLandline(t *testing.T) {
	// Create a validator instance and register our custom validation
	validate := validator.New()
	validate.RegisterValidation("taiwanlandline", ValidateTaiwanLandline)

	type TestStructPhone struct {
		Phone string `validate:"taiwanlandline"`
	}

	testStruct := TestStructPhone{Phone: "invalid-phone"}
	err := validate.Struct(testStruct)

	errors := ExtractValidationErrors(err)
	assert.Contains(t, errors, "phone")
	assert.Equal(t, "phone必須為有效的台灣市話號碼格式 (例: 02-12345678)", errors["phone"])
}
