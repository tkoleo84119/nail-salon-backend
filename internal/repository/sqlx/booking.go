package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
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

// UpdateMyBookingModel represents the database model for booking updates
type UpdateMyBookingModel struct {
	ID            int64  `db:"id"`
	StoreID       int64  `db:"store_id"`
	CustomerID    int64  `db:"customer_id"`
	StylistID     int64  `db:"stylist_id"`
	TimeSlotID    int64  `db:"time_slot_id"`
	IsChatEnabled bool   `db:"is_chat_enabled"`
	Note          string `db:"note"`
	Status        string `db:"status"`
}

// UpdateMyBooking updates a booking dynamically based on provided fields
func (r *BookingRepository) UpdateMyBooking(ctx context.Context, bookingID int64, customerID int64, req bookingModel.UpdateMyBookingRequest) (*UpdateMyBookingModel, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"booking_id":  bookingID,
		"customer_id": customerID,
	}

	// Dynamic field updates using type converters
	if req.StoreId != nil {
		storeID, err := utils.ParseID(*req.StoreId)
		if err != nil {
			return nil, fmt.Errorf("invalid store ID: %w", err)
		}
		setParts = append(setParts, "store_id = :store_id")
		args["store_id"] = storeID
	}

	if req.StylistId != nil {
		stylistID, err := utils.ParseID(*req.StylistId)
		if err != nil {
			return nil, fmt.Errorf("invalid stylist ID: %w", err)
		}
		setParts = append(setParts, "stylist_id = :stylist_id")
		args["stylist_id"] = stylistID
	}

	if req.TimeSlotId != nil {
		timeSlotID, err := utils.ParseID(*req.TimeSlotId)
		if err != nil {
			return nil, fmt.Errorf("invalid time slot ID: %w", err)
		}
		setParts = append(setParts, "time_slot_id = :time_slot_id")
		args["time_slot_id"] = timeSlotID
	}

	if req.IsChatEnabled != nil {
		setParts = append(setParts, "is_chat_enabled = :is_chat_enabled")
		args["is_chat_enabled"] = *req.IsChatEnabled
	}

	if req.Note != nil {
		setParts = append(setParts, "note = :note")
		args["note"] = utils.StringPtrToPgText(req.Note, true) // Empty as NULL
	}

	query := fmt.Sprintf(`
		UPDATE bookings SET %s
		WHERE id = :booking_id AND customer_id = :customer_id
		RETURNING id, store_id, customer_id, stylist_id, time_slot_id, is_chat_enabled, note, status
	`, strings.Join(setParts, ", "))

	var result UpdateMyBookingModel
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

// GetMyBookingsModel represents the database model for booking list
type GetMyBookingsModel struct {
	ID           int64  `db:"id"`
	StoreID      int64  `db:"store_id"`
	StoreName    string `db:"store_name"`
	StylistID    int64  `db:"stylist_id"`
	StylistName  string `db:"stylist_name"`
	Date         string `db:"date"`
	TimeSlotID   int64  `db:"time_slot_id"`
	StartTime    string `db:"start_time"`
	EndTime      string `db:"end_time"`
	Status       string `db:"status"`
}

// GetMyBookings retrieves customer bookings with optional filtering and pagination
func (r *BookingRepository) GetMyBookings(ctx context.Context, customerID int64, statuses []string, limit, offset int) ([]GetMyBookingsModel, int, error) {
	// Build WHERE conditions
	whereParts := []string{"b.customer_id = :customer_id"}
	args := map[string]interface{}{
		"customer_id": customerID,
		"limit":       limit,
		"offset":      offset,
	}

	// Add status filtering if provided
	if len(statuses) > 0 {
		placeholders := make([]string, len(statuses))
		for i, status := range statuses {
			placeholders[i] = fmt.Sprintf(":status_%d", i)
			args[fmt.Sprintf("status_%d", i)] = status
		}
		whereParts = append(whereParts, fmt.Sprintf("b.status IN (%s)", strings.Join(placeholders, ", ")))
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Query for bookings with joins
	query := fmt.Sprintf(`
		SELECT 
			b.id,
			b.store_id,
			s.name as store_name,
			b.stylist_id,
			st.name as stylist_name,
			TO_CHAR(sc.scheduled_date, 'YYYY-MM-DD') as date,
			b.time_slot_id,
			TO_CHAR(ts.start_time, 'HH24:MI') as start_time,
			TO_CHAR(ts.end_time, 'HH24:MI') as end_time,
			b.status
		FROM bookings b
		INNER JOIN stores s ON b.store_id = s.id
		INNER JOIN stylists st ON b.stylist_id = st.id
		INNER JOIN schedules sc ON b.stylist_id = sc.stylist_id AND b.time_slot_id = sc.time_slot_id
		INNER JOIN time_slots ts ON b.time_slot_id = ts.id
		WHERE %s
		ORDER BY sc.scheduled_date DESC, ts.start_time DESC
		LIMIT :limit OFFSET :offset
	`, whereClause)

	var bookings []GetMyBookingsModel
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("query bookings failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var booking GetMyBookingsModel
		if err := rows.StructScan(&booking); err != nil {
			return nil, 0, fmt.Errorf("scan booking failed: %w", err)
		}
		bookings = append(bookings, booking)
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM bookings b
		INNER JOIN stores s ON b.store_id = s.id
		INNER JOIN stylists st ON b.stylist_id = st.id
		INNER JOIN schedules sc ON b.stylist_id = sc.stylist_id AND b.time_slot_id = sc.time_slot_id
		INNER JOIN time_slots ts ON b.time_slot_id = ts.id
		WHERE %s
	`, whereClause)

	var total int
	countRow, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return nil, 0, fmt.Errorf("count bookings failed: %w", err)
	}
	defer countRow.Close()

	if countRow.Next() {
		if err := countRow.Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("scan count failed: %w", err)
		}
	}

	return bookings, total, nil
}
