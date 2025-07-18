package common

import "time"

// ParseTimeSlot parses a time slot string (HH:mm) into time.Time
func ParseTimeSlot(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// FormatTimeSlot formats a time.Time into a time slot string (HH:mm)
func FormatTimeSlot(t time.Time) string {
	return t.Format("15:04")
}
