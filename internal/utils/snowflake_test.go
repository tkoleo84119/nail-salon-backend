package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitSnowflake(t *testing.T) {
	snowflakeNode = nil
	once = sync.Once{}

	err := InitSnowflake(1)
	assert.NoError(t, err)
	assert.NotNil(t, snowflakeNode)

	// should not cause error, because it's a singleton
	err2 := InitSnowflake(2)
	assert.NoError(t, err2)
}

func TestGenerateID(t *testing.T) {
	// Ensure snowflake is initialized
	err := InitSnowflake(1)
	require.NoError(t, err)

	// Generate multiple IDs to test uniqueness
	ids := make(map[int64]bool)

	for i := 0; i < 1000; i++ {
		id := GenerateID()
		assert.True(t, id > 0, "ID should be positive")
		assert.False(t, ids[id], "ID should be unique")
		ids[id] = true
	}
}

func TestGenerateIDPanic(t *testing.T) {
	// Reset snowflake to nil
	snowflakeNode = nil

	// Should panic when snowflake is not initialized
	assert.Panics(t, func() {
		GenerateID()
	}, "Should panic when snowflake node is not initialized")

	// Reinitialize for other tests
	once = sync.Once{}
	InitSnowflake(1)
}
