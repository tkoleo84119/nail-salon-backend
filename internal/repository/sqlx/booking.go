package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type BookingRepositoryInterface interface {
	UpdateBookingTx(ctx context.Context, tx *sqlx.Tx, bookingID int64, params UpdateBookingParams) (int64, error)
	GetCustomerBookingByFilter(ctx context.Context, customerID int64, params GetCustomerBookingByFilterParams) ([]GetCustomerBookingByFilterItem, int, error)
	GetStoreBookingList(ctx context.Context, storeID int64, params GetStoreBookingListParams) ([]GetStoreBookingListModel, int, error)
	UpdateBookingByStaff(ctx context.Context, bookingID int64, storeID int64, timeSlotID *int64, isChatEnabled *bool, note *string) (*UpdateBookingByStaffModel, error)
	GetByID(ctx context.Context, bookingID int64) (*GetByIDRow, error)
	CancelBooking(ctx context.Context, tx *sqlx.Tx, bookingID int64, status string, cancelReason *string) (int64, error)
}

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

type UpdateBookingParams struct {
	StoreID       *int64
	StylistID     *int64
	TimeSlotID    *int64
	IsChatEnabled *bool
	Note          *string
}

// UpdateBooking updates a booking dynamically based on provided fields
func (r *BookingRepository) UpdateBookingTx(ctx context.Context, tx *sqlx.Tx, bookingID int64, params UpdateBookingParams) (int64, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{bookingID}

	if params.StoreID != nil {
		setParts = append(setParts, fmt.Sprintf("store_id = $%d", len(args)+1))
		args = append(args, *params.StoreID)
	}

	if params.StylistID != nil {
		setParts = append(setParts, fmt.Sprintf("stylist_id = $%d", len(args)+1))
		args = append(args, *params.StylistID)
	}

	if params.TimeSlotID != nil {
		setParts = append(setParts, fmt.Sprintf("time_slot_id = $%d", len(args)+1))
		args = append(args, *params.TimeSlotID)
	}

	if params.IsChatEnabled != nil {
		setParts = append(setParts, fmt.Sprintf("is_chat_enabled = $%d", len(args)+1))
		args = append(args, *params.IsChatEnabled)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	query := fmt.Sprintf(`
		UPDATE bookings
		SET %s
		WHERE id = $1
		RETURNING id
	`, strings.Join(setParts, ", "))

	var result int64
	err := tx.GetContext(ctx, &result, query, args...)
	if err != nil {
		return 0, fmt.Errorf("update booking failed: %w", err)
	}

	return result, nil
}

// ------------------------------------------------------------------------------------------------

type GetCustomerBookingByFilterParams struct {
	Limit    *int
	Offset   *int
	Sort     *[]string
	Statuses []string
}

type GetCustomerBookingByFilterItem struct {
	ID          int64       `db:"id"`
	StoreID     int64       `db:"store_id"`
	StoreName   string      `db:"store_name"`
	StylistID   int64       `db:"stylist_id"`
	StylistName string      `db:"stylist_name"`
	Date        pgtype.Date `db:"date"`
	TimeSlotID  int64       `db:"time_slot_id"`
	StartTime   pgtype.Time `db:"start_time"`
	EndTime     pgtype.Time `db:"end_time"`
	Status      string      `db:"status"`
}

func (r *BookingRepository) GetCustomerBookingByFilter(ctx context.Context, customerID int64, params GetCustomerBookingByFilterParams) (int, []GetCustomerBookingByFilterItem, error) {
	// Build WHERE conditions
	whereParts := []string{"b.customer_id = $1", "b.status != 'NO_SHOW'"}
	args := []interface{}{customerID}

	if len(params.Statuses) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("b.status = ANY($%d)", len(args)+1))
		args = append(args, params.Statuses)
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM bookings b
		INNER JOIN stores s ON b.store_id = s.id
		INNER JOIN stylists st ON b.stylist_id = st.id
		INNER JOIN time_slots ts ON b.time_slot_id = ts.id
		INNER JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
	`, whereClause)

	var total int
	row := r.db.QueryRowContext(ctx, countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return 0, nil, fmt.Errorf("count bookings failed: %w", err)
	}

	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	sort := utils.HandleSortByMap(map[string]string{
		"date":   "sch.work_date",
		"status": "b.status",
	}, "sch.work_date", "DESC", params.Sort)

	limitIndex := len(args) + 1
	offsetIndex := limitIndex + 1
	args = append(args, limit, offset)

	// Query for bookings with joins
	query := fmt.Sprintf(`
		SELECT
			b.id,
			b.store_id,
			s.name as store_name,
			b.stylist_id,
			st.name as stylist_name,
			sch.work_date as date,
			b.time_slot_id,
			ts.start_time as start_time,
			ts.end_time as end_time,
			b.status
		FROM bookings b
		INNER JOIN stores s ON b.store_id = s.id
		INNER JOIN stylists st ON b.stylist_id = st.id
		INNER JOIN time_slots ts ON b.time_slot_id = ts.id
		INNER JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var bookings []GetCustomerBookingByFilterItem
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("query bookings failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking GetCustomerBookingByFilterItem
		if err := rows.Scan(
			&booking.ID,
			&booking.StoreID,
			&booking.StoreName,
			&booking.StylistID,
			&booking.StylistName,
			&booking.Date,
			&booking.TimeSlotID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
		); err != nil {
			return 0, nil, fmt.Errorf("scan booking failed: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return total, bookings, nil
}

// ------------------------------------------------------------------------------------------------

// GetStoreBookingListParams represents the parameters for getting store booking list
type GetStoreBookingListParams struct {
	StylistID *int64
	StartDate *time.Time
	EndDate   *time.Time
	Limit     *int
	Offset    *int
}

// GetStoreBookingListModel represents the database model for admin booking list queries
type GetStoreBookingListModel struct {
	ID           int64       `db:"id"`
	StoreID      int64       `db:"store_id"`
	CustomerID   int64       `db:"customer_id"`
	CustomerName string      `db:"customer_name"`
	StylistID    int64       `db:"stylist_id"`
	StylistName  pgtype.Text `db:"stylist_name"`
	TimeSlotID   int64       `db:"time_slot_id"`
	StartTime    pgtype.Time `db:"start_time"`
	EndTime      pgtype.Time `db:"end_time"`
	WorkDate     pgtype.Date `db:"work_date"`
	Status       string      `db:"status"`
}

// GetStoreBookingList retrieves bookings for a specific store with dynamic filtering and pagination
func (r *BookingRepository) GetStoreBookingList(ctx context.Context, storeID int64, params GetStoreBookingListParams) ([]GetStoreBookingListModel, int, error) {
	// Set default pagination values
	limit := 20
	offset := 0
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}

	// Build WHERE clause parts
	whereParts := []string{"b.store_id = :store_id"}
	args := map[string]interface{}{
		"store_id": storeID,
		"limit":    limit,
		"offset":   offset,
	}

	// Add stylist filter
	if params.StylistID != nil {
		whereParts = append(whereParts, "b.stylist_id = :stylist_id")
		args["stylist_id"] = *params.StylistID
	}

	// Add start date filter
	if params.StartDate != nil {
		whereParts = append(whereParts, "sch.work_date >= :start_date")
		args["start_date"] = *params.StartDate
	}

	// Add end date filter
	if params.EndDate != nil {
		whereParts = append(whereParts, "sch.work_date <= :end_date")
		args["end_date"] = *params.EndDate
	}

	// Build WHERE clause
	whereClause := strings.Join(whereParts, " AND ")

	// Count query for total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM bookings b
		JOIN time_slots ts ON b.time_slot_id = ts.id
		JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
	`, whereClause)

	var total int
	rows, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT
			b.id,
			b.store_id,
			b.customer_id,
			c.name as customer_name,
			b.stylist_id,
			st.name as stylist_name,
			b.time_slot_id,
			ts.start_time,
			ts.end_time,
			sch.work_date,
			b.status
		FROM bookings b
		JOIN customers c ON b.customer_id = c.id
		JOIN stylists st ON b.stylist_id = st.id
		JOIN time_slots ts ON b.time_slot_id = ts.id
		JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
		ORDER BY sch.work_date ASC, ts.start_time ASC
		LIMIT :limit OFFSET :offset
	`, whereClause)

	var results []GetStoreBookingListModel
	rows, err = r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetStoreBookingListModel
		if err := rows.StructScan(&item); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, total, nil
}

// UpdateBookingByStaffModel represents the database model for booking updates by staff
type UpdateBookingByStaffModel struct {
	ID            int64       `db:"id"`
	StoreID       int64       `db:"store_id"`
	CustomerID    int64       `db:"customer_id"`
	StylistID     int64       `db:"stylist_id"`
	TimeSlotID    int64       `db:"time_slot_id"`
	IsChatEnabled pgtype.Bool `db:"is_chat_enabled"`
	Note          pgtype.Text `db:"note"`
	Status        string      `db:"status"`
}

// UpdateBookingByStaff updates a booking dynamically based on provided fields for staff
func (r *BookingRepository) UpdateBookingByStaff(ctx context.Context, bookingID int64, storeID int64, timeSlotID *int64, isChatEnabled *bool, note *string) (*UpdateBookingByStaffModel, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"booking_id": bookingID,
		"store_id":   storeID,
	}

	// Dynamic field updates
	if timeSlotID != nil {
		setParts = append(setParts, "time_slot_id = :time_slot_id")
		args["time_slot_id"] = *timeSlotID
	}

	if isChatEnabled != nil {
		setParts = append(setParts, "is_chat_enabled = :is_chat_enabled")
		args["is_chat_enabled"] = *isChatEnabled
	}

	if note != nil {
		setParts = append(setParts, "note = :note")
		args["note"] = utils.StringPtrToPgText(note, false)
	}

	query := fmt.Sprintf(`
		UPDATE bookings SET %s
		WHERE id = :booking_id AND store_id = :store_id AND status = 'SCHEDULED'
		RETURNING id, store_id, customer_id, stylist_id, time_slot_id, is_chat_enabled, note, status
	`, strings.Join(setParts, ", "))

	var result UpdateBookingByStaffModel
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &result, nil
}

type GetByIDRow struct {
	ID            int64              `db:"id"`
	StoreID       int64              `db:"store_id"`
	StoreName     string             `db:"store_name"`
	CustomerID    int64              `db:"customer_id"`
	StylistID     int64              `db:"stylist_id"`
	StylistName   pgtype.Text        `db:"stylist_name"`
	TimeSlotID    int64              `db:"time_slot_id"`
	StartTime     pgtype.Time        `db:"start_time"`
	EndTime       pgtype.Time        `db:"end_time"`
	WorkDate      pgtype.Date        `db:"work_date"`
	IsChatEnabled pgtype.Bool        `db:"is_chat_enabled"`
	Note          pgtype.Text        `db:"note"`
	Status        string             `db:"status"`
	CreatedAt     pgtype.Timestamptz `db:"created_at"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at"`
}

// GetByID retrieves a booking by ID
func (r *BookingRepository) GetByID(ctx context.Context, bookingID int64) (*GetByIDRow, error) {
	query := `
		SELECT
			b.id,
			b.store_id,
			s.name as store_name,
			b.customer_id,
			b.stylist_id,
			st.name as stylist_name,
			b.time_slot_id,
			ts.start_time,
			ts.end_time,
			sch.work_date,
			b.is_chat_enabled,
			b.note,
			b.status,
			b.created_at,
			b.updated_at
		FROM bookings b
		JOIN stores s ON b.store_id = s.id
		JOIN stylists st ON b.stylist_id = st.id
		JOIN time_slots ts ON b.time_slot_id = ts.id
		JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE b.id = $1
	`

	var result GetByIDRow
	err := r.db.GetContext(ctx, &result, query, bookingID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelBooking cancels a booking with transaction support
func (r *BookingRepository) CancelBooking(ctx context.Context, tx *sqlx.Tx, bookingID int64, status string, cancelReason *string) (int64, error) {
	cancelParams := map[string]interface{}{
		"booking_id":    bookingID,
		"status":        status,
		"cancel_reason": utils.StringPtrToPgText(cancelReason, false),
	}

	cancelQuery := `
		UPDATE bookings
		SET status = :status, cancel_reason = :cancel_reason, updated_at = NOW()
		WHERE id = :booking_id
		RETURNING id
	`

	var result int64
	rows, err := tx.NamedQuery(cancelQuery, cancelParams)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel booking: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("no rows returned after cancelling booking")
	}

	if err := rows.StructScan(&result); err != nil {
		return 0, fmt.Errorf("failed to scan booking ID: %w", err)
	}

	return result, nil
}
