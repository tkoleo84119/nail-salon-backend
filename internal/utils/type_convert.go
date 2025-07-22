package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringToText converts a string pointer to pgtype.Text
func StringToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// TextToString converts pgtype.Text to string, returns empty string if invalid
func TextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// NumericToFloat64 converts pgtype.Numeric to float64
func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	f, err := n.Float64Value()
	if err != nil {
		return 0
	}

	return f.Float64
}

// Float64ToNumeric converts float64 to pgtype.Numeric
func Float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	err := n.Scan(f)
	if err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

// BoolPtrToPgBool converts *bool to pgtype.Bool
func BoolPtrToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// TimeToPgTimez converts time.Time to pgtype.Timestamptz
func TimeToPgTimez(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// TimeToString converts pgtype.Time to string format HH:MM
func TimeToString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	// Convert microseconds to time
	totalMicros := t.Microseconds
	hours := totalMicros / (60 * 60 * 1000000)
	minutes := (totalMicros % (60 * 60 * 1000000)) / (60 * 1000000)

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// DateToString converts pgtype.Date to string format YYYY-MM-DD
func DateToString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}

	return d.Time.Format("2006-01-02")
}

// StringToTime converts date string to time.Time
func StringToTime(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// StringToPgDate converts date string to pgtype.Date
func StringToPgDate(s string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}, err
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

// TimeToPgDate converts time.Time to pgtype.Date
func TimeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}
