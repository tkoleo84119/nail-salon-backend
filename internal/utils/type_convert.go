package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------- Pgtype conversion to Go type functions ----------------------------------

func PgTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func PgNumericToFloat64(n pgtype.Numeric) (float64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("invalid numeric value")
	}

	f, err := n.Float64Value()
	if err != nil {
		return 0, err
	}
	return f.Float64, nil
}

func PgNumericToInt64(n pgtype.Numeric) (int64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("invalid numeric value")
	}

	f, err := n.Int64Value()
	if err != nil {
		return 0, err
	}
	return f.Int64, nil
}

func PgTimeToTimeString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	totalMicros := t.Microseconds
	hours := totalMicros / (60 * 60 * 1000000)
	minutes := (totalMicros % (60 * 60 * 1000000)) / (60 * 1000000)

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func PgTimeToTime(t pgtype.Time) (time.Time, error) {
	if !t.Valid {
		return time.Time{}, fmt.Errorf("invalid time")
	}

	d := time.Duration(t.Microseconds) * time.Microsecond
	return time.Unix(0, 0).UTC().Add(d), nil
}

func PgDateToDateString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func PgInt8ToIDString(id pgtype.Int8) string {
	if !id.Valid {
		return ""
	}
	return FormatID(id.Int64)
}

func PgInt4ToInt32(value pgtype.Int4) int32 {
	if !value.Valid {
		return 0
	}
	return value.Int32
}

func PgInt4ToInt32Ptr(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

func PgBoolToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func PgTimestamptzToTimeString(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

// ---------------------------------- Go type conversion to Pgtype functions ----------------------------------

func StringPtrToPgText(s *string, emptyAsNull bool) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}

	if emptyAsNull && *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func BoolPtrToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func Float64PtrToPgNumeric(f *float64) (pgtype.Numeric, error) {
	if f == nil {
		return pgtype.Numeric{Valid: false}, nil
	}

	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%f", *f))
	if err != nil {
		return pgtype.Numeric{Valid: false}, err
	}
	return n, nil
}

func Int64PtrToPgNumeric(i *int64) (pgtype.Numeric, error) {
	if i == nil {
		return pgtype.Numeric{Valid: false}, nil
	}

	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%d", *i))
	if err != nil {
		return pgtype.Numeric{Valid: false}, err
	}
	return n, nil
}

func TimePtrToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func TimePtrToPgTime(t *time.Time) pgtype.Time {
	if t == nil {
		return pgtype.Time{Valid: false}
	}

	if t.IsZero() {
		return pgtype.Time{Valid: false}
	}

	totalMicros := int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000
	totalMicros += int64(t.Nanosecond()) / 1000 // Add microsecond precision
	return pgtype.Time{Microseconds: totalMicros, Valid: true}
}

func TimePtrToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}

	if t.IsZero() {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{Time: *t, Valid: true}
}

func DateStringToPgDate(s string) (pgtype.Date, error) {
	t, err := DateStringToTime(s)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	return pgtype.Date{Time: t, Valid: true}, nil
}

func TimeStringToPgTime(s string) (pgtype.Time, error) {
	t, err := TimeStringToTime(s)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("invalid time format, expected HH:MM: %w", err)
	}

	totalMicros := int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000
	return pgtype.Time{Microseconds: totalMicros, Valid: true}, nil
}

func Int64PtrToPgInt8(id *int64) pgtype.Int8 {
	if id == nil {
		return pgtype.Int8{Valid: false}
	}
	if *id == 0 {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *id, Valid: true}
}

func Int32PtrToPgInt4(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

// ---------------------------------- String conversion to Go type functions ----------------------------------

func DateStringToTime(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func DateStringToTimeInLoc(s string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func TimeStringToTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
