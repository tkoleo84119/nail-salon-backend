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

