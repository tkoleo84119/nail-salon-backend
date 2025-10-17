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

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllCustomerBookingsByFilterParams struct {
	Limit  *int
	Offset *int
	Sort   *[]string
	Status []string
}

type GetAllCustomerBookingsByFilterItem struct {
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

// GetAllCustomerBookingsByFilter retrieves all bookings for a customer with dynamic filtering and pagination
func (r *BookingRepository) GetAllCustomerBookingsByFilter(ctx context.Context, customerID int64, params GetAllCustomerBookingsByFilterParams) (int, []GetAllCustomerBookingsByFilterItem, error) {
	// WHERE conditions
	whereParts := []string{"b.customer_id = $1", "b.status != 'NO_SHOW'"}
	args := []interface{}{customerID}

	if len(params.Status) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("b.status = ANY($%d)", len(args)+1))
		args = append(args, params.Status)
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM bookings b
		WHERE %s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("count bookings failed: %w", err)
	}
	if total == 0 {
		return 0, []GetAllCustomerBookingsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"sch.work_date DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"date":   "sch.work_date",
		"status": "b.status",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

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

	var bookings []GetAllCustomerBookingsByFilterItem
	if err := r.db.SelectContext(ctx, &bookings, query, args...); err != nil {
		return 0, nil, fmt.Errorf("query bookings failed: %w", err)
	}

	return total, bookings, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStoreBookingsByFilterParams struct {
	StylistID  *int64
	CustomerID *int64
	StartDate  *time.Time
	EndDate    *time.Time
	Status     *string
	Limit      *int
	Offset     *int
	Sort       *[]string
}

type GetAllStoreBookingsByFilterItem struct {
	ID               int64              `db:"id"`
	StoreID          int64              `db:"store_id"`
	CustomerID       int64              `db:"customer_id"`
	CustomerName     string             `db:"customer_name"`
	CustomerLineName string             `db:"customer_line_name"`
	StylistID        int64              `db:"stylist_id"`
	StylistName      pgtype.Text        `db:"stylist_name"`
	TimeSlotID       int64              `db:"time_slot_id"`
	StartTime        pgtype.Time        `db:"start_time"`
	EndTime          pgtype.Time        `db:"end_time"`
	WorkDate         pgtype.Date        `db:"work_date"`
	ActualDuration   pgtype.Int4        `db:"actual_duration"`
	Status           string             `db:"status"`
	CreatedAt        pgtype.Timestamptz `db:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at"`
}

// GetAllStoreBookingsByFilter retrieves bookings for a specific store with dynamic filtering and pagination
func (r *BookingRepository) GetAllStoreBookingsByFilter(ctx context.Context, storeID int64, params GetAllStoreBookingsByFilterParams) (int, []GetAllStoreBookingsByFilterItem, error) {
	// WHERE conditions
	whereConditions := []string{"b.store_id = $1"}
	args := []interface{}{storeID}

	if params.StylistID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.stylist_id = $%d", len(args)+1))
		args = append(args, *params.StylistID)
	}

	if params.CustomerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.customer_id = $%d", len(args)+1))
		args = append(args, *params.CustomerID)
	}

	if params.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sch.work_date >= $%d", len(args)+1))
		args = append(args, *params.StartDate)
	}

	if params.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sch.work_date <= $%d", len(args)+1))
		args = append(args, *params.EndDate)
	}

	if params.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("b.status = $%d", len(args)+1))
		args = append(args, *params.Status)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM bookings b
		JOIN time_slots ts ON b.time_slot_id = ts.id
		JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("count bookings failed: %w", err)
	}
	if total == 0 {
		return 0, []GetAllStoreBookingsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"sch.work_date DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"date":      "sch.work_date",
		"status":    "b.status",
		"customer":  "c.name",
		"stylist":   "st.name",
		"startTime": "ts.start_time",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIdx := len(args) - 1
	offsetIdx := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT
			b.id,
			b.store_id,
			b.customer_id,
			c.name as customer_name,
			COALESCE(c.line_name, '') as customer_line_name,
			b.stylist_id,
			st.name as stylist_name,
			b.time_slot_id,
			ts.start_time,
			ts.end_time,
			sch.work_date,
			b.actual_duration,
			b.status,
			b.created_at,
			b.updated_at
		FROM bookings b
		JOIN customers c ON b.customer_id = c.id
		JOIN stylists st ON b.stylist_id = st.id
		JOIN time_slots ts ON b.time_slot_id = ts.id
		JOIN schedules sch ON ts.schedule_id = sch.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIdx, offsetIdx)

	var results []GetAllStoreBookingsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("query bookings failed: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateBookingTxParams struct {
	StoreID       *int64
	StylistID     *int64
	TimeSlotID    *int64
	IsChatEnabled *bool
	Note          *string
	StoreNote     *string
}

// UpdateBooking updates a booking dynamically based on provided fields
func (r *BookingRepository) UpdateBookingTx(ctx context.Context, tx *sqlx.Tx, bookingID int64, params UpdateBookingTxParams) (int64, error) {
	// SET conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

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

	if params.StoreNote != nil {
		setParts = append(setParts, fmt.Sprintf("store_note = $%d", len(args)+1))
		args = append(args, *params.StoreNote)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return 0, fmt.Errorf("no fields to update")
	}

	args = append(args, bookingID)
	setClause := strings.Join(setParts, ", ")

	// Data query
	query := fmt.Sprintf(`
		UPDATE bookings
		SET %s
		WHERE id = $%d
		RETURNING id
	`, setClause, len(args))

	var result int64
	if err := tx.GetContext(ctx, &result, query, args...); err != nil {
		return 0, fmt.Errorf("update booking failed: %w", err)
	}

	return result, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// CancelBookingTx cancels a booking with transaction support
func (r *BookingRepository) CancelBookingTx(ctx context.Context, tx *sqlx.Tx, bookingID int64, status string, cancelReason *string) (int64, error) {
	// Data query
	query := `
		UPDATE bookings
		SET status = $1,
			cancel_reason = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING id
  `

	args := []interface{}{
		status,
		utils.StringPtrToPgText(cancelReason, true),
		bookingID,
	}

	var result int64
	if err := tx.GetContext(ctx, &result, query, args...); err != nil {
		return 0, fmt.Errorf("failed to cancel booking: %w", err)
	}

	return result, nil
}
