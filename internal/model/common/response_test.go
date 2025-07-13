package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	data := map[string]interface{}{
		"id":   123,
		"name": "test",
	}

	response := SuccessResponse(data)

	assert.Empty(t, response.Message)
	assert.Equal(t, data, response.Data)
	assert.Nil(t, response.Errors)
}

func TestErrorResponse(t *testing.T) {
	message := "Something went wrong"
	errors := map[string]string{
		"field1": "error message 1",
		"field2": "error message 2",
	}

	response := ErrorResponse(message, errors)

	assert.Equal(t, message, response.Message)
	assert.Nil(t, response.Data)
	assert.Equal(t, errors, response.Errors)
}

func TestValidationErrorResponse(t *testing.T) {
	errors := map[string]string{
		"username": "帳號為必填項目",
		"password": "密碼長度至少需要8個字元",
	}

	response := ValidationErrorResponse(errors)

	assert.Equal(t, "輸入驗證失敗", response.Message)
	assert.Nil(t, response.Data)
	assert.Equal(t, errors, response.Errors)
}