package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// StringPtrToPgText converts string pointer to pgtype.Text for nullable string fields.
// This is the unified function for handling optional string fields with configurable empty string behavior.
//
// Parameters:
//   - s: Pointer to string (can be nil for NULL values)
//   - emptyAsNull: Whether to treat empty strings as NULL (true) or as valid empty strings (false)
//
// Returns:
//   - pgtype.Text with Valid=false if s is nil, or if emptyAsNull=true and s points to empty string
//   - pgtype.Text with Valid=true and the string value otherwise
//
// Examples:
//
//	// Basic usage - nil becomes NULL, empty string stays as empty string
//	var name *string = nil
//	pgText := utils.StringPtrToPgText(name, false)          // Returns {Valid: false}
//
//	emptyName := ""
//	pgText = utils.StringPtrToPgText(&emptyName, false)     // Returns {String: "", Valid: true}
//
//	realName := "Alice"
//	pgText = utils.StringPtrToPgText(&realName, false)      // Returns {String: "Alice", Valid: true}
//
//	// Empty as NULL behavior - both nil and empty string become NULL
//	var address *string = nil
//	pgText = utils.StringPtrToPgText(address, true)         // Returns {Valid: false}
//
//	emptyAddr := ""
//	pgText = utils.StringPtrToPgText(&emptyAddr, true)      // Returns {Valid: false}
//
//	realAddr := "123 Main St"
//	pgText = utils.StringPtrToPgText(&realAddr, true)       // Returns {String: "123 Main St", Valid: true}
//
// Usage in SQLC queries:
//
//	// For optional fields where empty strings are valid
//	params := CreateUserParams{
//	    Name:    req.Name,                                    // Required field
//	    Bio:     utils.StringPtrToPgText(req.Bio, false),    // Optional, empty allowed
//	}
//
//	// For optional fields where empty strings should be NULL
//	params := UpdateStoreParams{
//	    Address: utils.StringPtrToPgText(req.Address, true), // Optional, empty becomes NULL
//	    Phone:   utils.StringPtrToPgText(req.Phone, true),   // Optional, empty becomes NULL
//	}
func StringPtrToPgText(s *string, emptyAsNull bool) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	if emptyAsNull && *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// PgTextToString converts pgtype.Text to string, returning empty string for invalid/NULL values.
// This is used when reading nullable string fields from database results.
//
// Parameters:
//   - t: pgtype.Text from database query result
//
// Returns:
//   - Empty string if t.Valid is false (NULL in database)
//   - The actual string value if t.Valid is true
//
// Example:
//
//	var result struct {
//	    Name    string      `db:"name"`
//	    Address pgtype.Text `db:"address"`
//	}
//	// After scanning from database...
//	addressStr := utils.PgTextToString(result.Address)  // "" if NULL, actual value otherwise
//
// Usage in response building:
//
//	response := &UserResponse{
//	    Name:    result.Name,
//	    Address: utils.PgTextToString(result.Address),
//	}
func PgTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// PgNumericToFloat64 converts pgtype.Numeric to float64 for numeric calculations.
// Returns 0 for invalid/NULL values or conversion errors.
//
// Parameters:
//   - n: pgtype.Numeric from database query result
//
// Returns:
//   - 0 if n.Valid is false (NULL in database) or conversion fails
//   - The float64 value if conversion succeeds
//
// Example:
//
//	var result struct {
//	    Price pgtype.Numeric `db:"price"`
//	}
//	// After scanning from database...
//	priceFloat := utils.PgNumericToFloat64(result.Price)  // 0 if NULL, actual value otherwise
//
// Usage in response building:
//
//	response := &ProductResponse{
//	    ID:    result.ID,
//	    Price: utils.PgNumericToFloat64(result.Price),
//	}
func PgNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}

	f, err := n.Float64Value()
	if err != nil {
		return 0
	}

	return f.Float64
}

// Float64ToPgNumeric converts float64 to pgtype.Numeric with error handling.
// This is used when inserting/updating numeric values in database operations.
//
// Parameters:
//   - f: float64 value to convert
//
// Returns:
//   - pgtype.Numeric with Valid=true and the converted value
//   - error if conversion fails
//
// Example:
//
//	price := 99.99
//	pgPrice, err := utils.Float64ToPgNumeric(price)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to convert price: %w", err)
//	}
//
// Usage in SQLC queries:
//
//	priceNumeric, err := utils.Float64ToPgNumeric(req.Price)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert price", err)
//	}
//	params := CreateProductParams{
//	    Name:  req.Name,
//	    Price: priceNumeric,
//	}
func Float64ToPgNumeric(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(f)
	if err != nil {
		return pgtype.Numeric{Valid: false}, err
	}
	return n, nil
}

// Int64ToPgNumeric converts int64 to pgtype.Numeric with error handling.
// This is used when inserting/updating integer values as numeric in database operations.
//
// Parameters:
//   - i: int64 value to convert
//
// Returns:
//   - pgtype.Numeric with Valid=true and the converted value
//   - error if conversion fails
//
// Example:
//
//	price := int64(1000) // Price in cents
//	pgPrice, err := utils.Int64ToPgNumeric(price)
//	if err != nil {
//	    return nil, fmt.Errorf("failed to convert price: %w", err)
//	}
//
// Usage in SQLC queries:
//
//	priceNumeric, err := utils.Int64ToPgNumeric(req.Price)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to convert price", err)
//	}
//	params := CreateServiceParams{
//	    Name:  req.Name,
//	    Price: priceNumeric,
//	}
func Int64ToPgNumeric(i int64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%d", i))
	if err != nil {
		return pgtype.Numeric{Valid: false}, err
	}
	return n, nil
}

// BoolPtrToPgBool converts a boolean pointer to pgtype.Bool for nullable boolean fields.
// This is used for optional boolean fields in database operations where NULL values are allowed.
//
// Parameters:
//   - b: Pointer to bool (can be nil for NULL values)
//
// Returns:
//   - pgtype.Bool with Valid=false if b is nil, otherwise Valid=true with the boolean value
//
// Example:
//
//	var isActive *bool = nil
//	pgBool := utils.BoolPtrToPgBool(isActive)          // Returns {Valid: false}
//
//	enabled := true
//	pgBool = utils.BoolPtrToPgBool(&enabled)           // Returns {Bool: true, Valid: true}
//
// Usage in SQLC queries:
//
//	params := CreateServiceParams{
//	    Name:      req.Name,
//	    IsVisible: utils.BoolPtrToPgBool(&req.IsVisible),  // Required field
//	    IsAddon:   utils.BoolPtrToPgBool(req.IsAddon),     // Optional field
//	}
func BoolPtrToPgBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// Float64PtrToPgNumeric converts a float64 pointer to pgtype.Numeric for nullable float fields.
// This is used for optional float fields in database operations where NULL values are allowed.
//
// Parameters:
//   - f: Pointer to float64 (can be nil for NULL values)
//
// Returns:
//   - pgtype.Numeric with Valid=false if f is nil, otherwise Valid=true with the converted value
//
// Example:
//
//	var discountRate *float64 = nil
//	pgNumeric := utils.Float64PtrToPgNumeric(discountRate)  // Returns {Valid: false}
//
//	rate := 0.8
//	pgNumeric = utils.Float64PtrToPgNumeric(&rate)          // Returns {Valid: true, value: 0.8}
//
// Usage in SQLC queries:
//
//	params := CreateCouponParams{
//	    Name:          req.Name,
//	    DiscountRate:  utils.Float64PtrToPgNumeric(req.DiscountRate),   // Optional field
//	    DiscountAmount: utils.Int64PtrToPgNumeric(req.DiscountAmount),  // Optional field
//	}
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

// Int64PtrToPgNumeric converts an int64 pointer to pgtype.Numeric for nullable int fields.
// This is used for optional int fields in database operations where NULL values are allowed.
//
// Parameters:
//   - i: Pointer to int64 (can be nil for NULL values)
//
// Returns:
//   - pgtype.Numeric with Valid=false if i is nil, otherwise Valid=true with the converted value
//
// Example:
//
//	var discountAmount *int64 = nil
//	pgNumeric := utils.Int64PtrToPgNumeric(discountAmount)  // Returns {Valid: false}
//
//	amount := int64(100)
//	pgNumeric = utils.Int64PtrToPgNumeric(&amount)          // Returns {Valid: true, value: 100}
//
// Usage in SQLC queries:
//
//	params := CreateCouponParams{
//	    Name:          req.Name,
//	    DiscountRate:  utils.Float64PtrToPgNumeric(req.DiscountRate),   // Optional field
//	    DiscountAmount: utils.Int64PtrToPgNumeric(req.DiscountAmount),  // Optional field
//	}
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

// TimeToPgTimestamptz converts time.Time to pgtype.Timestamptz for timestamp fields.
// This is used for created_at, updated_at, and other timestamp fields in database operations.
//
// Parameters:
//   - t: time.Time value to convert
//
// Returns:
//   - pgtype.Timestamptz with Valid=true and the timestamp value
//
// Example:
//
//	now := time.Now()
//	pgTimestamp := utils.TimeToPgTimestamptz(now)
//
// Usage in SQLC queries:
//
//	params := CreateUserParams{
//	    Name:      req.Name,
//	    CreatedAt: utils.TimeToPgTimestamptz(time.Now()),
//	    UpdatedAt: utils.TimeToPgTimestamptz(time.Now()),
//	}
func TimeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{Time: t, Valid: true}
}

// TimePtrToPgTimestamptz converts a pointer to time.Time to pgtype.Timestamptz for nullable timestamp fields.
// This is used for optional timestamp fields in database operations where NULL values are allowed.
//
// Parameters:
//   - t: Pointer to time.Time (can be nil for NULL values)
//
// Returns:
//   - pgtype.Timestamptz with Valid=false if t is nil, otherwise Valid=true with the timestamp value
//
// Example:
//
//	var createdAt *time.Time = nil
//	pgTimestamp := utils.TimePtrToPgTimestamptz(createdAt)  // Returns {Valid: false}
//
//	now := time.Now()
//	pgTimestamp = utils.TimePtrToPgTimestamptz(&now)       // Returns {Time: now, Valid: true}
//
// Usage in SQLC queries:
//
//	params := CreateUserParams{
//	    Name:      req.Name,
//	    CreatedAt: utils.TimePtrToPgTimestamptz(req.CreatedAt),  // Optional field
//	    UpdatedAt: utils.TimePtrToPgTimestamptz(req.UpdatedAt),  // Optional field
//	}
func TimePtrToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// TimeToPgTime converts time.Time to pgtype.Time for time-of-day fields.
// This extracts the time portion (HH:MM:SS) and converts to PostgreSQL time type with microsecond precision.
//
// Parameters:
//   - t: time.Time value to convert (only time portion is used)
//
// Returns:
//   - pgtype.Time with Valid=true and the time value in microseconds
//
// Example:
//
//	timeValue, _ := time.Parse("15:04", "09:30")
//	pgTime := utils.TimeToPgTime(timeValue)           // Converts to 09:30:00
//
// Usage in SQLC queries:
//
//	startTime, _ := utils.TimeStringToTime(req.StartTime)
//	params := CreateTimeSlotParams{
//	    StartTime: utils.TimeToPgTime(startTime),
//	    EndTime:   utils.TimeToPgTime(endTime),
//	}
func TimeToPgTime(t time.Time) pgtype.Time {
	totalMicros := int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000
	totalMicros += int64(t.Nanosecond()) / 1000 // Add microsecond precision
	return pgtype.Time{Microseconds: totalMicros, Valid: true}
}

// PgTimeToTimeString converts pgtype.Time to string in HH:MM format.
// This is used when building API responses that need time in string format.
//
// Parameters:
//   - t: pgtype.Time from database query result
//
// Returns:
//   - Empty string if t.Valid is false (NULL in database)
//   - Time string in "HH:MM" format (e.g., "09:30", "14:45")
//
// Example:
//
//	var result struct {
//	    StartTime pgtype.Time `db:"start_time"`
//	}
//	// After scanning from database...
//	timeStr := utils.PgTimeToTimeString(result.StartTime)  // "09:30"
//
// Usage in response building:
//
//	response := &TimeSlotResponse{
//	    ID:        utils.FormatID(result.ID),
//	    StartTime: utils.PgTimeToTimeString(result.StartTime),
//	    EndTime:   utils.PgTimeToTimeString(result.EndTime),
//	}
func PgTimeToTimeString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	totalMicros := t.Microseconds
	hours := totalMicros / (60 * 60 * 1000000)
	minutes := (totalMicros % (60 * 60 * 1000000)) / (60 * 1000000)

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// PgTimeToTime converts pgtype.Time to time.Time for time calculations.
// This is used when you need to perform time operations or comparisons.
//
// Parameters:
//   - t: pgtype.Time from database query result
//
// Returns:
//   - time.Time with the time portion set (date will be 0001-01-01)
//   - error if the time is invalid
//
// Example:
//
//	var result struct {
//	    StartTime pgtype.Time `db:"start_time"`
//	}
//	// After scanning from database...
//	startTime, err := utils.PgTimeToTime(result.StartTime)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert time", err)
//	}
//
// Usage in business logic:
//
//	startTime, err := utils.PgTimeToTime(timeSlot.StartTime)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert time", err)
//	}
//	endTime, err := utils.PgTimeToTime(timeSlot.EndTime)
//	if err !=	 nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to convert time", err)
//	}
//
//	// Validate time range
//	if !endTime.After(startTime) {
//	    return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
//	}
func PgTimeToTime(t pgtype.Time) (time.Time, error) {
	if !t.Valid {
		return time.Time{}, fmt.Errorf("invalid time")
	}

	d := time.Duration(t.Microseconds) * time.Microsecond
	return time.Unix(0, 0).UTC().Add(d), nil
}

// PgDateToDateString converts pgtype.Date to string in YYYY-MM-DD format.
// This is used when building API responses that need date in string format.
//
// Parameters:
//   - d: pgtype.Date from database query result
//
// Returns:
//   - Empty string if d.Valid is false (NULL in database)
//   - Date string in "YYYY-MM-DD" format (e.g., "2023-12-25")
//
// Example:
//
//	var result struct {
//	    WorkDate pgtype.Date `db:"work_date"`
//	}
//	// After scanning from database...
//	dateStr := utils.PgDateToDateString(result.WorkDate)  // "2023-12-25"
//
// Usage in response building:
//
//	response := &ScheduleResponse{
//	    ID:       utils.FormatID(result.ID),
//	    WorkDate: utils.PgDateToDateString(result.WorkDate),
//	}
func PgDateToDateString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// DateStringToTime converts date string (YYYY-MM-DD) to time.Time.
// This is used for parsing date strings from API requests or other sources.
//
// Parameters:
//   - s: Date string in "YYYY-MM-DD" format
//
// Returns:
//   - time.Time with the parsed date
//   - error if the string format is invalid
//
// Example:
//
//	dateStr := "2023-12-25"
//	date, err := utils.DateStringToTime(dateStr)
//	if err != nil {
//	    return fmt.Errorf("invalid date: %w", err)
//	}
//
// Usage in request validation:
//
//	workDate, err := utils.DateStringToTime(req.WorkDate)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid work date format", err)
//	}
func DateStringToTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// DateStringToPgDate converts date string (YYYY-MM-DD) to pgtype.Date with error handling.
// This is used when inserting/updating date values in database operations.
//
// Parameters:
//   - s: Date string in "YYYY-MM-DD" format
//
// Returns:
//   - pgtype.Date with Valid=true and the converted date
//   - error if the string format is invalid
//
// Example:
//
//	workDate, err := utils.DateStringToPgDate("2023-12-25")
//	if err != nil {
//	    return nil, fmt.Errorf("invalid date format: %w", err)
//	}
//
// Usage in SQLC queries:
//
//	workDate, err := utils.DateStringToPgDate(req.WorkDate)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid work date format", err)
//	}
//	params := CreateScheduleParams{
//	    WorkDate: workDate,
//	}
func DateStringToPgDate(s string) (pgtype.Date, error) {
	t, err := DateStringToTime(s)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

// TimeStringToTime converts time string (HH:MM) to time.Time for time calculations.
// This is used for parsing time strings from API requests or other sources.
//
// Parameters:
//   - s: Time string in "HH:MM" format (e.g., "09:30", "14:45")
//
// Returns:
//   - time.Time with the parsed time (date will be 0001-01-01)
//   - error if the string format is invalid
//
// Example:
//
//	timeStr := "09:30"
//	timeValue, err := utils.TimeStringToTime(timeStr)
//	if err != nil {
//	    return fmt.Errorf("invalid time: %w", err)
//	}
//
// Usage in request validation:
//
//	startTime, err := utils.TimeStringToTime(req.StartTime)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid start time format", err)
//	}
//
//	endTime, err := utils.TimeStringToTime(req.EndTime)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid end time format", err)
//	}
//
//	// Validate time range
//	if !endTime.After(startTime) {
//	    return nil, errorCodes.NewServiceErrorWithCode(errorCodes.TimeSlotEndBeforeStart)
//	}
func TimeStringToTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

// TimeStringToPgTime converts time string (HH:MM) to pgtype.Time with error handling.
// This is used when inserting/updating time values in database operations.
//
// Parameters:
//   - s: Time string in "HH:MM" format (e.g., "09:30", "14:45")
//
// Returns:
//   - pgtype.Time with Valid=true and the converted time
//   - error if the string format is invalid
//
// Example:
//
//	pgTime, err := utils.TimeStringToPgTime("09:30")
//	if err != nil {
//	    return nil, fmt.Errorf("invalid time format: %w", err)
//	}
//
// Usage in SQLC queries:
//
//	startTime, err := utils.TimeStringToPgTime(req.StartTime)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid start time format", err)
//	}
//	params := CreateTimeSlotParams{
//	    StartTime: startTime,
//	}
func TimeStringToPgTime(s string) (pgtype.Time, error) {
	t, err := TimeStringToTime(s)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("invalid time format, expected HH:MM: %w", err)
	}

	totalMicros := int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000
	return pgtype.Time{Microseconds: totalMicros, Valid: true}, nil
}

// TimeToPgDate converts time.Time to pgtype.Date for date-only fields.
// This extracts the date portion and converts to PostgreSQL date type.
//
// Parameters:
//   - t: time.Time value to convert (only date portion is used)
//
// Returns:
//   - pgtype.Date with Valid=true and the date value
//
// Example:
//
//	now := time.Now()
//	pgDate := utils.TimeToPgDate(now)                  // Converts to date-only
//
// Usage in SQLC queries:
//
//	params := CreateScheduleParams{
//	    WorkDate: utils.TimeToPgDate(time.Now()),
//	}
func TimeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// Nullable ID conversions - for handling optional foreign key relationships

// Int64ToPgInt8 converts int64 to pgtype.Int8 for nullable ID fields.
// This is used when setting foreign key references that are known to be valid.
//
// Parameters:
//   - id: int64 value to convert (typically from utils.ParseID())
//
// Returns:
//   - pgtype.Int8 with Valid=true and the ID value
//
// Example:
//
//	staffUserID, _ := utils.ParseID(staffContext.UserID)
//	updaterID := utils.Int64ToPgInt8(staffUserID)
//
// Usage in SQLC queries:
//
//	params := CreateTimeSlotTemplateParams{
//	    Name:    req.Name,
//	    Updater: utils.Int64ToPgInt8(staffUserID),  // Foreign key reference
//	}
func Int64ToPgInt8(id int64) pgtype.Int8 {
	return pgtype.Int8{Int64: id, Valid: true}
}

// Int64PtrToPgInt8 converts *int64 to pgtype.Int8 for optional ID fields.
// This is used when foreign key references might be NULL.
//
// Parameters:
//   - id: Pointer to int64 (can be nil for NULL values)
//
// Returns:
//   - pgtype.Int8 with Valid=false if id is nil, otherwise Valid=true with the ID value
//
// Example:
//
//	var couponID *int64 = nil
//	pgCouponID := utils.Int64PtrToPgInt8(couponID)     // Returns {Valid: false}
//
//	payerID := int64(123)
//	pgPayerID := utils.Int64PtrToPgInt8(&payerID)      // Returns {Int64: 123, Valid: true}
//
// Usage in SQLC queries:
//
//	params := CreateCheckoutParams{
//	    BookingID:    bookingID,                       // Required foreign key
//	    CouponID:     utils.Int64PtrToPgInt8(req.CouponID),  // Optional foreign key
//	    CheckoutUser: utils.Int64PtrToPgInt8(req.ProcessedBy), // Optional foreign key
//	}
func Int64PtrToPgInt8(id *int64) pgtype.Int8 {
	if id == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *id, Valid: true}
}

// Int32ToPgInt4 converts int32 to pgtype.Int4 for nullable integer fields.
// This is used when setting integer values that are known to be valid.
//
// Parameters:
//   - value: int32 value to convert
//
// Returns:
//   - pgtype.Int4 with Valid=true and the integer value
//
// Example:
//
//	actualDuration := int32(45)  // minutes
//	pgDuration := utils.Int32ToPgInt4(actualDuration)
//
// Usage in SQLC queries:
//
//	params := UpdateBookingParams{
//	    ActualDuration: utils.Int32ToPgInt4(req.ActualDuration),
//	}
func Int32ToPgInt4(value int32) pgtype.Int4 {
	return pgtype.Int4{Int32: value, Valid: true}
}

// Int32PtrToPgInt4 converts *int32 to pgtype.Int4 for optional integer fields.
// This is used when integer values might be NULL.
//
// Parameters:
//   - value: Pointer to int32 (can be nil for NULL values)
//
// Returns:
//   - pgtype.Int4 with Valid=false if value is nil, otherwise Valid=true with the integer value
//
// Example:
//
//	var safetyStock *int32 = nil
//	pgStock := utils.Int32PtrToPgInt4(safetyStock)     // Returns {Valid: false}
//
//	minStock := int32(10)
//	pgStock = utils.Int32PtrToPgInt4(&minStock)        // Returns {Int32: 10, Valid: true}
//
// Usage in SQLC queries:
//
//	params := CreateProductParams{
//	    Name:        req.Name,                         // Required field
//	    SafetyStock: utils.Int32PtrToPgInt4(req.SafetyStock), // Optional field
//	}
func Int32PtrToPgInt4(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

// ParseIDToPgInt8 parses string ID to pgtype.Int8 with validation.
// This combines ID parsing and pgtype conversion in one step.
//
// Parameters:
//   - idStr: ID string to parse (empty string is treated as NULL)
//
// Returns:
//   - pgtype.Int8 with Valid=false for empty string, Valid=true for valid ID
//   - error if ID parsing fails
//
// Example:
//
//	couponID, err := utils.ParseIDToPgInt8(req.CouponID)
//	if err != nil {
//	    return nil, fmt.Errorf("invalid coupon ID: %w", err)
//	}
//
// Usage in request processing:
//
//	couponID, err := utils.ParseIDToPgInt8(req.CouponID)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid coupon ID", err)
//	}
//	params := CreateCheckoutParams{
//	    CouponID: couponID,
//	}
func ParseIDToPgInt8(idStr string) (pgtype.Int8, error) {
	if idStr == "" {
		return pgtype.Int8{Valid: false}, nil
	}

	id, err := ParseID(idStr)
	if err != nil {
		return pgtype.Int8{Valid: false}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

// ParseIDPtrToPgInt8 parses optional string ID to pgtype.Int8.
// This handles pointer to string for optional ID fields.
//
// Parameters:
//   - idStr: Pointer to ID string (can be nil or point to empty string for NULL)
//
// Returns:
//   - pgtype.Int8 with Valid=false for nil or empty string, Valid=true for valid ID
//   - error if ID parsing fails
//
// Example:
//
//	var couponID *string = nil
//	pgCouponID, err := utils.ParseIDPtrToPgInt8(couponID)  // Returns {Valid: false}, nil
//
//	payerID := "123"
//	pgPayerID, err := utils.ParseIDPtrToPgInt8(&payerID)   // Returns {Int64: 123, Valid: true}, nil
//
// Usage in request processing:
//
//	payerID, err := utils.ParseIDPtrToPgInt8(req.PayerID)
//	if err != nil {
//	    return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid payer ID", err)
//	}
//	params := CreateExpenseParams{
//	    PayerID: payerID,
//	}
func ParseIDPtrToPgInt8(idStr *string) (pgtype.Int8, error) {
	if idStr == nil || *idStr == "" {
		return pgtype.Int8{Valid: false}, nil
	}

	return ParseIDToPgInt8(*idStr)
}

// PgInt8ToIDString converts pgtype.Int8 to string for API responses.
// This is used when building API responses that include optional ID fields.
//
// Parameters:
//   - id: pgtype.Int8 from database query result
//
// Returns:
//   - Empty string if id.Valid is false (NULL in database)
//   - Formatted ID string if id.Valid is true
//
// Example:
//
//	var result struct {
//	    CouponID pgtype.Int8 `db:"coupon_id"`
//	}
//	// After scanning from database...
//	couponIDStr := utils.PgInt8ToIDString(result.CouponID)  // "" if NULL, "123" otherwise
//
// Usage in response building:
//
//	response := &CheckoutResponse{
//	    ID:       utils.FormatID(result.ID),
//	    CouponID: utils.PgInt8ToIDString(result.CouponID),  // Optional foreign key
//	}
func PgInt8ToIDString(id pgtype.Int8) string {
	if !id.Valid {
		return ""
	}
	return FormatID(id.Int64)
}

// PgInt8ToIDStringPtr converts pgtype.Int8 to *string for optional API fields.
// This is used when API responses need to distinguish between empty and NULL values.
//
// Parameters:
//   - id: pgtype.Int8 from database query result
//
// Returns:
//   - nil if id.Valid is false (NULL in database)
//   - Pointer to formatted ID string if id.Valid is true
//
// Example:
//
//	var result struct {
//	    CouponID pgtype.Int8 `db:"coupon_id"`
//	}
//	// After scanning from database...
//	couponIDPtr := utils.PgInt8ToIDStringPtr(result.CouponID)  // nil if NULL, &"123" otherwise
//
// Usage in response building:
//
//	response := &CheckoutResponse{
//	    ID:       utils.FormatID(result.ID),
//	    CouponID: utils.PgInt8ToIDStringPtr(result.CouponID),  // Optional field as pointer
//	}
func PgInt8ToIDStringPtr(id pgtype.Int8) *string {
	if !id.Valid {
		return nil
	}
	idStr := FormatID(id.Int64)
	return &idStr
}

// PgInt4ToInt32 converts pgtype.Int4 to int32 for API responses.
// Returns 0 for invalid/NULL values.
//
// Parameters:
//   - value: pgtype.Int4 from database query result
//
// Returns:
//   - 0 if value.Valid is false (NULL in database)
//   - The int32 value if value.Valid is true
//
// Example:
//
//	var result struct {
//	    ActualDuration pgtype.Int4 `db:"actual_duration"`
//	}
//	// After scanning from database...
//	duration := utils.PgInt4ToInt32(result.ActualDuration)  // 0 if NULL, actual value otherwise
//
// Usage in response building:
//
//	response := &BookingResponse{
//	    ID:             utils.FormatID(result.ID),
//	    ActualDuration: utils.PgInt4ToInt32(result.ActualDuration),
//	}
func PgInt4ToInt32(value pgtype.Int4) int32 {
	if !value.Valid {
		return 0
	}
	return value.Int32
}

// PgInt4ToInt32Ptr converts pgtype.Int4 to *int32 for optional API fields.
// This preserves the distinction between NULL and 0 values.
//
// Parameters:
//   - value: pgtype.Int4 from database query result
//
// Returns:
//   - nil if value.Valid is false (NULL in database)
//   - Pointer to int32 value if value.Valid is true
//
// Example:
//
//	var result struct {
//	    SafetyStock pgtype.Int4 `db:"safety_stock"`
//	}
//	// After scanning from database...
//	stockPtr := utils.PgInt4ToInt32Ptr(result.SafetyStock)  // nil if NULL, &10 otherwise
//
// Usage in response building:
//
//	response := &ProductResponse{
//	    ID:          utils.FormatID(result.ID),
//	    SafetyStock: utils.PgInt4ToInt32Ptr(result.SafetyStock),  // Optional field as pointer
//	}
func PgInt4ToInt32Ptr(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

// PgBoolToBool converts pgtype.Bool to bool for API responses.
// Returns false for invalid/NULL values.
//
// Parameters:
//   - b: pgtype.Bool from database query result
//
// Returns:
//   - false if b.Valid is false (NULL in database)
//   - The bool value if b.Valid is true
//
// Example:
//
//	var result struct {
//	    IsAvailable pgtype.Bool `db:"is_available"`
//	}
//	// After scanning from database...
//	available := utils.PgBoolToBool(result.IsAvailable)  // false if NULL, actual value otherwise
//
// Usage in response building:
//
//	response := &TimeSlotResponse{
//	    ID:          utils.FormatID(result.ID),
//	    IsAvailable: utils.PgBoolToBool(result.IsAvailable),
//	}
func PgBoolToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

// BoolToPgBool converts bool to pgtype.Bool for database operations.
// This is used when setting boolean values that are known to be valid.
//
// Parameters:
//   - b: bool value to convert
//
// Returns:
//   - pgtype.Bool with Valid=true and the boolean value
//
// Example:
//
//	isEnabled := true
//	pgBool := utils.BoolToPgBool(isEnabled)
//
// Usage in SQLC queries:
//
//	params := CreateBookingParams{
//	    IsChatEnabled: utils.BoolToPgBool(req.IsChatEnabled),
//	}
func BoolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

// PgTimestamptzToTimeString converts pgtype.Timestamptz to string in ISO 8601 format.
// This is used when building API responses that need timestamp in string format.
//
// Parameters:
//   - t: pgtype.Timestamptz from database query result
//
// Returns:
//   - Empty string if t.Valid is false (NULL in database)
//   - Time string in "YYYY-MM-DDTHH:MM:SS+08:00" format (e.g., "2025-01-01T00:00:00+08:00")
//
// Example:
//
//	var result struct {
//	    CreatedAt pgtype.Timestamptz `db:"created_at"`
//	}
//	// After scanning from database...
//	createdAtStr := utils.PgTimestamptzToTimeString(result.CreatedAt)  // "" if NULL, "2025-01-01T00:00:00+08:00" otherwise
//
// Usage in response building:
//
//		response := &BookingResponse{
//		    ID:          utils.FormatID(result.ID),
//		    CreatedAt:   utils.PgTimestamptzToTimeString(result.CreatedAt),
//		}
//	}
func PgTimestamptzToTimeString(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}
