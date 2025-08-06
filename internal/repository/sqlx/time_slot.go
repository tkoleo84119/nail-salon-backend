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
	BatchCreateTimeSlotsTx(ctx context.Context, tx *sqlx.Tx, params []BatchCreateTimeSlotsTxParams) error
	UpdateTimeSlot(ctx context.Context, timeSlotID int64, params UpdateTimeSlotParams) (UpdateTimeSlotResponse, error)
}

type TimeSlotRepository struct {
	db *sqlx.DB
}

func NewTimeSlotRepository(db *sqlx.DB) *TimeSlotRepository {
	return &TimeSlotRepository{db: db}
}

type BatchCreateTimeSlotsTxParams struct {
	ID         int64
	ScheduleID int64
	StartTime  pgtype.Time
	EndTime    pgtype.Time
}

func (r *TimeSlotRepository) BatchCreateTimeSlotsTx(ctx context.Context, tx *sqlx.Tx, params []BatchCreateTimeSlotsTxParams) error {
	const batchSize = 1000

	var (
		sb   strings.Builder
		args []interface{}
	)

	for i := 0; i < len(params); i += batchSize {
		end := i + batchSize
		if end > len(params) {
			end = len(params)
		}

		sb.Reset()
		args = args[:0]

		sb.WriteString("INSERT INTO time_slots (id, schedule_id, start_time, end_time) VALUES ")

		param := 1
		for j, v := range params[i:end] {
			sb.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d)", param, param+1, param+2, param+3))
			if j < end-i-1 {
				sb.WriteByte(',')
			}
			args = append(args, v.ID, v.ScheduleID, v.StartTime, v.EndTime)
			param += 4
		}

		if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	return nil
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
