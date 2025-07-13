package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRefreshToken(t *testing.T) {
	// Generate multiple tokens to test uniqueness
	tokens := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		token, err := GenerateRefreshToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, 64, len(token)) // 32 bytes = 64 hex characters
		
		// Check uniqueness
		assert.False(t, tokens[token], "Token should be unique")
		tokens[token] = true
		
		// Check that it's valid hex
		for _, char := range token {
			assert.True(t, 
				(char >= '0' && char <= '9') || (char >= 'a' && char <= 'f'),
				"Token should only contain hex characters")
		}
	}
}