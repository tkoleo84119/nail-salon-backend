package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TimeSlotRepositoryInterface defines the interface for time slot repository
type TimeSlotRepositoryInterface interface {
	UpdateTimeSlot(ctx context.Context, timeSlotID int64, params UpdateTimeSlotParams) (UpdateTimeSlotResponse, error)
	UpdateTimeSlotAvailabilityTx(ctx context.Context, tx *sqlx.Tx, timeSlotID int64, isAvailable bool) error
}

type TimeSlotRepository struct {
	db *sqlx.DB
}

func NewTimeSlotRepository(db *sqlx.DB) *TimeSlotRepository {
	return &TimeSlotRepository{
		db: db,
	}
}

type UpdateTimeSlotParams struct {
	StartTime   *string
	EndTime     *string
	IsAvailable *bool
}

type UpdateTimeSlotResponse struct {
	ID          int64
	ScheduleID  int64
	StartTime   pgtype.Time
	EndTime     pgtype.Time
	IsAvailable pgtype.Bool
}

// UpdateTimeSlot updates time slot with dynamic fields
func (r *TimeSlotRepository) UpdateTimeSlot(ctx context.Context, timeSlotID int64, params UpdateTimeSlotParams) (UpdateTimeSlotResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{timeSlotID}

	if params.StartTime != nil {
		setParts = append(setParts, fmt.Sprintf("start_time = $%d", len(args)+1))
		startTime, err := utils.TimeStringToPgTime(*params.StartTime)
		if err != nil {
			return UpdateTimeSlotResponse{}, fmt.Errorf("failed to convert start time: %w", err)
		}
		args = append(args, startTime)
	}

	if params.EndTime != nil {
		setParts = append(setParts, fmt.Sprintf("end_time = $%d", len(args)+1))
		endTime, err := utils.TimeStringToPgTime(*params.EndTime)
		if err != nil {
			return UpdateTimeSlotResponse{}, fmt.Errorf("failed to convert end time: %w", err)
		}
		args = append(args, endTime)
	}

	if params.IsAvailable != nil {
		setParts = append(setParts, fmt.Sprintf("is_available = $%d", len(args)+1))
		args = append(args, params.IsAvailable)
	}

	query := fmt.Sprintf(`
		UPDATE time_slots
		SET %s
		WHERE id = $1
		RETURNING
			id,
			schedule_id,
			start_time,
			end_time,
			is_available
	`, strings.Join(setParts, ", "))

	var result UpdateTimeSlotResponse
	rows := r.db.QueryRowContext(ctx, query, args...)
	if err := rows.Scan(
		&result.ID,
		&result.ScheduleID,
		&result.StartTime,
		&result.EndTime,
		&result.IsAvailable,
	); err != nil {
		return UpdateTimeSlotResponse{}, fmt.Errorf("failed to scan result: %w", err)
	}

	return result, nil
}

// UpdateTimeSlotAvailability updates the availability status of a time slot
func (r *TimeSlotRepository) UpdateTimeSlotAvailability(ctx context.Context, timeSlotID int64, isAvailable bool) error {
	query := `
		UPDATE time_slots
		SET is_available = :is_available, updated_at = NOW()
		WHERE id = :id
	`

	args := map[string]interface{}{
		"id":           timeSlotID,
		"is_available": isAvailable,
	}

	_, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update time slot availability: %w", err)
	}

	return nil
}

// UpdateTimeSlotAvailabilityTx updates the availability status of a time slot with transaction support
func (r *TimeSlotRepository) UpdateTimeSlotAvailabilityTx(ctx context.Context, tx *sqlx.Tx, timeSlotID int64, isAvailable bool) error {
	query := `
		UPDATE time_slots
		SET is_available = :is_available, updated_at = NOW()
		WHERE id = :id
	`

	args := map[string]interface{}{
		"id":           timeSlotID,
		"is_available": isAvailable,
	}

	_, err := tx.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update time slot availability: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (r *TimeSlotRepository) UpdateTimeSlotAvailabilityByBookingIDTx(ctx context.Context, tx *sqlx.Tx, timeSlotID int64, isAvailable bool) error {
	query := `
		UPDATE time_slots
		SET is_available = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, isAvailable, timeSlotID)
	if err != nil {
		return fmt.Errorf("failed to update time slot availability: %w", err)
	}

	return nil
}
