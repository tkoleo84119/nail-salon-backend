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
