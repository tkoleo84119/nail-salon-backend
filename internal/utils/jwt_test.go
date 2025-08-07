package utils

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
)

func TestGenerateJWT(t *testing.T) {
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key",
		ExpiryHours: 1,
	}

	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:    "valid JWT generation",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "another valid user",
			userID:  456,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(jwtConfig, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGenerateJWTWithEmptySecret(t *testing.T) {
	jwtConfig := config.JWTConfig{
		Secret:      "", // Empty secret
		ExpiryHours: 1,
	}

	token, err := GenerateJWT(jwtConfig, 123)
	// JWT library allows empty secrets, so this should not error
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key",
		ExpiryHours: 1,
	}

	userID := int64(123)

	token, err := GenerateJWT(jwtConfig, userID)
	require.NoError(t, err)

	claims, err := ValidateJWT(jwtConfig, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, strconv.FormatInt(userID, 10), claims.UserID)
}

func TestValidateJWTInvalid(t *testing.T) {
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key",
		ExpiryHours: 1,
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid token format",
			token: "invalid-token",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "malformed JWT",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(jwtConfig, tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestValidateJWTExpired(t *testing.T) {
	jwtConfig := config.JWTConfig{
		Secret:      "test-secret-key",
		ExpiryHours: 0, // Set to 0 hours to make it expire immediately
	}

	token, err := GenerateJWT(jwtConfig, 123)
	require.NoError(t, err)

	time.Sleep(time.Second)

	claims, err := ValidateJWT(jwtConfig, token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
