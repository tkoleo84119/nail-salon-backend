package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeSlot(t *testing.T) {
	tests := []struct {
		name        string
		timeStr     string
		expected    time.Time
		expectError bool
	}{
		{
			name:        "Valid time format",
			timeStr:     "09:30",
			expected:    time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Valid time format with zero minutes",
			timeStr:     "14:00",
			expected:    time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Valid time format midnight",
			timeStr:     "00:00",
			expected:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Valid time format end of day",
			timeStr:     "23:59",
			expected:    time.Date(2025, 1, 1, 23, 59, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Invalid format - missing colon",
			timeStr:     "0930",
			expectError: true,
		},
		{
			name:        "Invalid format - wrong separator",
			timeStr:     "09.30",
			expectError: true,
		},
		{
			name:        "Invalid format - 12-hour format",
			timeStr:     "9:30 AM",
			expectError: true,
		},
		{
			name:        "Valid format - single digit hour",
			timeStr:     "9:30",
			expected:    time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Invalid format - invalid hour",
			timeStr:     "25:30",
			expectError: true,
		},
		{
			name:        "Invalid format - invalid minute",
			timeStr:     "09:61",
			expectError: true,
		},
		{
			name:        "Empty string",
			timeStr:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimeSlot(tt.timeStr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Hour(), result.Hour())
				assert.Equal(t, tt.expected.Minute(), result.Minute())
			}
		})
	}
}

func TestFormatTimeSlot(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Morning time",
			time:     time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC),
			expected: "09:30",
		},
		{
			name:     "Afternoon time",
			time:     time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC),
			expected: "14:00",
		},
		{
			name:     "Midnight",
			time:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "00:00",
		},
		{
			name:     "End of day",
			time:     time.Date(2025, 1, 1, 23, 59, 0, 0, time.UTC),
			expected: "23:59",
		},
		{
			name:     "Single digit minute",
			time:     time.Date(2025, 1, 1, 10, 5, 0, 0, time.UTC),
			expected: "10:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeSlot(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}
