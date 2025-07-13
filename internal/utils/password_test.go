package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "long password exceeds bcrypt limit",
			password: "this-is-a-very-long-password-that-should-still-work-fine-with-bcrypt-hashing-algorithm-but-exceeds-72-bytes-limit",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash) // Hash should be different from original
				assert.True(t, len(hash) >= 60)       // bcrypt hashes are at least 60 chars
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		want           bool
	}{
		{
			name:           "correct password",
			password:       password,
			hashedPassword: hash,
			want:           true,
		},
		{
			name:           "incorrect password",
			password:       "wrongpassword",
			hashedPassword: hash,
			want:           false,
		},
		{
			name:           "empty password",
			password:       "",
			hashedPassword: hash,
			want:           false,
		},
		{
			name:           "invalid hash",
			password:       password,
			hashedPassword: "invalid-hash",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hashedPassword)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	passwords := []string{
		"simple",
		"complex-P@ssw0rd!",
		"emoji-password-🔐",
		"very-long-password-with-many-characters-to-test-edge-cases",
	}

	for _, password := range passwords {
		t.Run("password: "+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Check correct password
			assert.True(t, CheckPassword(password, hash))

			// Check incorrect password
			assert.False(t, CheckPassword(password+"wrong", hash))
		})
	}
}