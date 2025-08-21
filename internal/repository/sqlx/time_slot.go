package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type TimeSlotRepository struct {
	db *sqlx.DB
}

func NewTimeSlotRepository(db *sqlx.DB) *TimeSlotRepository {
	return &TimeSlotRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateTimeSlotParams struct {
	StartTime   *string
	EndTime     *string
	IsAvailable *bool
}

type UpdateTimeSlotResponse struct {
	ID          int64       `db:"id"`
	ScheduleID  int64       `db:"schedule_id"`
	StartTime   pgtype.Time `db:"start_time"`
	EndTime     pgtype.Time `db:"end_time"`
	IsAvailable pgtype.Bool `db:"is_available"`
}

// UpdateTimeSlot updates time slot with dynamic fields
func (r *TimeSlotRepository) UpdateTimeSlot(ctx context.Context, timeSlotID int64, params UpdateTimeSlotParams) (UpdateTimeSlotResponse, error) {
	// Set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.StartTime != nil {
		setParts = append(setParts, fmt.Sprintf("start_time = $%d", len(args)+1))
		startTime, err := utils.TimeStringToPgTime(*params.StartTime)
		if err != nil {
			return UpdateTimeSlotResponse{}, fmt.Errorf("convert start time failed: %w", err)
		}
		args = append(args, startTime)
	}

	if params.EndTime != nil {
		setParts = append(setParts, fmt.Sprintf("end_time = $%d", len(args)+1))
		endTime, err := utils.TimeStringToPgTime(*params.EndTime)
		if err != nil {
			return UpdateTimeSlotResponse{}, fmt.Errorf("convert end time failed: %w", err)
		}
		args = append(args, endTime)
	}

	if params.IsAvailable != nil {
		setParts = append(setParts, fmt.Sprintf("is_available = $%d", len(args)+1))
		args = append(args, params.IsAvailable)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateTimeSlotResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, timeSlotID)

	// Data query
	query := fmt.Sprintf(`
		UPDATE time_slots
		SET %s
		WHERE id = $%d
		RETURNING
			id,
			schedule_id,
			start_time,
			end_time,
			is_available
	`, strings.Join(setParts, ", "), len(args))

	var result UpdateTimeSlotResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateTimeSlotResponse{}, fmt.Errorf("update time slot failed: %w", err)
	}

	return result, nil
}

// ---------------------------------------------------------------------------------------------------------------------

// UpdateTimeSlotAvailabilityTx updates the availability status of a time slot with transaction support
func (r *TimeSlotRepository) UpdateTimeSlotAvailabilityTx(ctx context.Context, tx *sqlx.Tx, timeSlotID int64, isAvailable bool) error {
	query := `
		UPDATE time_slots
		SET is_available = $1, updated_at = NOW()
		WHERE id = $2
	`

	args := []interface{}{isAvailable, timeSlotID}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update time slot availability failed: %w", err)
	}

	return nil
}
