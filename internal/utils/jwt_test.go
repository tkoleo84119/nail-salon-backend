package utils

import (
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
		name      string
		userID    int64
		username  string
		role      string
		storeList []Store
		wantErr   bool
	}{
		{
			name:     "valid JWT generation",
			userID:   123,
			username: "testuser",
			role:     "ADMIN",
			storeList: []Store{
				{ID: 1, Name: "Store 1"},
				{ID: 2, Name: "Store 2"},
				{ID: 3, Name: "Store 3"},
			},
			wantErr: false,
		},
		{
			name:      "empty store IDs",
			userID:    456,
			username:  "manager",
			role:      "MANAGER",
			storeList: []Store{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(jwtConfig, tt.userID, tt.username, tt.role, tt.storeList)

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

	token, err := GenerateJWT(jwtConfig, 123, "test", "ADMIN", []Store{
		{ID: 1, Name: "Store 1"},
	})
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
	username := "testuser"
	role := "ADMIN"
	storeList := []Store{
		{ID: 1, Name: "Store 1"},
		{ID: 2, Name: "Store 2"},
		{ID: 3, Name: "Store 3"},
	}

	token, err := GenerateJWT(jwtConfig, userID, username, role, storeList)
	require.NoError(t, err)

	claims, err := ValidateJWT(jwtConfig, token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, storeList, claims.StoreList)
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

	token, err := GenerateJWT(jwtConfig, 123, "test", "ADMIN", []Store{
		{ID: 1, Name: "Store 1"},
	})
	require.NoError(t, err)

	time.Sleep(time.Second)

	claims, err := ValidateJWT(jwtConfig, token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
