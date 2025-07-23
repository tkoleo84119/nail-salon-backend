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

func Int64ToPgNumeric(i int64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(i)
	if err != nil {
		return pgtype.Numeric{Valid: false}, err
	}
	return n, nil
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

// TimeToPgTime converts time.Time to pgtype.Time
func TimeToPgTime(t time.Time) pgtype.Time {
	return pgtype.Time{Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000, Valid: true}
}

// PgTimeToStringTime converts pgtype.Time to string format HH:MM
func PgTimeToStringTime(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	// Convert microseconds to time
	totalMicros := t.Microseconds
	hours := totalMicros / (60 * 60 * 1000000)
	minutes := (totalMicros % (60 * 60 * 1000000)) / (60 * 1000000)

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func PgTimeToTime(t pgtype.Time) time.Time {
	return time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(t.Microseconds) * time.Microsecond)
}

// PgDateToString converts pgtype.Date to string format YYYY-MM-DD
func PgDateToString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}

	return d.Time.Format("2006-01-02")
}

// StringDateToTime converts date string (YYYY-MM-DD) to time.Time
func StringDateToTime(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// StringDateToPgDate converts date string (YYYY-MM-DD) to pgtype.Date
func StringDateToPgDate(s string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}, err
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

// StringTimeToTime converts time string (HH:MM) to time.Time
func StringTimeToTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// StringTimeToPgTime converts time string (HH:MM) to pgtype.Time
func StringTimeToPgTime(s string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return pgtype.Time{}, err
	}
	return pgtype.Time{Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000, Valid: true}, nil
}

// TimeToPgDate converts time.Time to pgtype.Date
func TimeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}
