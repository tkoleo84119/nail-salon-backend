package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TimeSlotRepositoryInterface defines the interface for time slot repository
type TimeSlotRepositoryInterface interface {
	BatchCreateTimeSlotsTx(ctx context.Context, tx *sqlx.Tx, params []BatchCreateTimeSlotsTxParams) error
	UpdateTimeSlot(ctx context.Context, timeSlotID int64, req adminScheduleModel.UpdateTimeSlotRequest) (*adminScheduleModel.UpdateTimeSlotResponse, error)
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

// UpdateTimeSlot updates time slot with dynamic fields
func (r *TimeSlotRepository) UpdateTimeSlot(ctx context.Context, timeSlotID int64, req adminScheduleModel.UpdateTimeSlotRequest) (*adminScheduleModel.UpdateTimeSlotResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": timeSlotID,
	}

	if req.StartTime != nil {
		setParts = append(setParts, "start_time = :start_time")
		// Convert time string to pgtype.Time
		startTime, err := utils.TimeStringToTime(*req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start time format: %w", err)
		}
		args["start_time"] = startTime
	}

	if req.EndTime != nil {
		setParts = append(setParts, "end_time = :end_time")
		// Convert time string to pgtype.Time
		endTime, err := utils.TimeStringToTime(*req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end time format: %w", err)
		}
		args["end_time"] = endTime
	}

	if req.IsAvailable != nil {
		setParts = append(setParts, "is_available = :is_available")
		args["is_available"] = *req.IsAvailable
	}

	query := fmt.Sprintf(`
		UPDATE time_slots
		SET %s
		WHERE id = :id
		RETURNING
			id,
			schedule_id,
			start_time,
			end_time,
			is_available,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "))

	var result dbgen.TimeSlot
	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows returned from update")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	response := &adminScheduleModel.UpdateTimeSlotResponse{
		ID:          utils.FormatID(result.ID),
		ScheduleID:  utils.FormatID(result.ScheduleID),
		StartTime:   utils.PgTimeToTimeString(result.StartTime),
		EndTime:     utils.PgTimeToTimeString(result.EndTime),
		IsAvailable: result.IsAvailable.Bool,
	}

	return response, nil
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
