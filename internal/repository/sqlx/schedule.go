package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

// ScheduleRepositoryInterface defines the interface for schedule repository
type ScheduleRepositoryInterface interface {
	BatchCreateSchedulesTx(ctx context.Context, tx *sqlx.Tx, params []BatchCreateSchedulesTxParams) error
	GetStoreScheduleByDateRange(ctx context.Context, storeID int64, startDate time.Time, endDate time.Time, params GetStoreScheduleByDateRangeParams) ([]GetStoreScheduleByDateRangeItem, error)
	GetScheduleByID(ctx context.Context, scheduleID int64) ([]GetScheduleByIDItem, error)
	CheckScheduleExists(ctx context.Context, storeID int64, stylistID int64, workDate time.Time) (bool, error)
	UpdateSchedule(ctx context.Context, scheduleID int64, params UpdateScheduleParams) (UpdateScheduleResponse, error)
}

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

type BatchCreateSchedulesTxParams struct {
	ID        int64
	StoreID   int64
	StylistID int64
	WorkDate  pgtype.Date
	Note      pgtype.Text
}

func (r *ScheduleRepository) BatchCreateSchedulesTx(ctx context.Context, tx *sqlx.Tx, params []BatchCreateSchedulesTxParams) error {
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

		sb.WriteString(
			"INSERT INTO schedules (id, store_id, stylist_id, work_date, note) VALUES ",
		)

		param := 1
		for j, v := range params[i:end] {
			sb.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", param, param+1, param+2, param+3, param+4))
			if j < end-i-1 {
				sb.WriteByte(',')
			}
			args = append(args, v.ID, v.StoreID, v.StylistID, v.WorkDate, v.Note)
			param += 5
		}

		if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}
	return nil
}

type GetStoreScheduleByDateRangeParams struct {
	StylistID   *[]int64
	IsAvailable *bool
}

type GetStoreScheduleByDateRangeItem struct {
	ID          int64       `db:"id"`
	StylistID   int64       `db:"stylist_id"`
	WorkDate    pgtype.Date `db:"work_date"`
	Note        pgtype.Text `db:"note"`
	StylistName pgtype.Text `db:"stylist_name"`
	TimeSlotID  int64       `db:"time_slot_id"`
	StartTime   pgtype.Time `db:"start_time"`
	EndTime     pgtype.Time `db:"end_time"`
	IsAvailable pgtype.Bool `db:"is_available"`
}

// GetStoreScheduleList retrieves schedules for a specific store using step-by-step queries
func (r *ScheduleRepository) GetStoreScheduleByDateRange(ctx context.Context, storeID int64, startDate time.Time, endDate time.Time, params GetStoreScheduleByDateRangeParams) ([]GetStoreScheduleByDateRangeItem, error) {
	scheduleWhereParts := []string{
		"s.store_id = $1",
		"s.work_date BETWEEN $2 AND $3",
	}
	scheduleArgs := []interface{}{storeID, startDate, endDate}

	// Add stylist filter if provided
	if params.StylistID != nil && len(*params.StylistID) > 0 {
		scheduleWhereParts = append(scheduleWhereParts, "s.stylist_id = ANY($4)")
		scheduleArgs = append(scheduleArgs, *params.StylistID)
	}

	whereClause := strings.Join(scheduleWhereParts, " AND ")

	scheduleQuery := fmt.Sprintf(`
		SELECT
			s.id,
			s.stylist_id,
			s.work_date,
			COALESCE(s.note, '') as note,
			st.name as stylist_name,
			ts.id as time_slot_id,
			ts.start_time,
			ts.end_time,
			ts.is_available
		FROM schedules s
		LEFT JOIN stylists st ON s.stylist_id = st.id
		LEFT JOIN time_slots ts ON s.id = ts.schedule_id
		WHERE %s
		ORDER BY s.work_date ASC, ts.start_time ASC
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, scheduleQuery, scheduleArgs...)
	if err != nil {
		return []GetStoreScheduleByDateRangeItem{}, err
	}
	defer rows.Close()

	response := []GetStoreScheduleByDateRangeItem{}
	for rows.Next() {
		var item GetStoreScheduleByDateRangeItem
		err := rows.Scan(
			&item.ID,
			&item.StylistID,
			&item.WorkDate,
			&item.Note,
			&item.StylistName,
			&item.TimeSlotID,
			&item.StartTime,
			&item.EndTime,
			&item.IsAvailable,
		)
		if err != nil {
			return []GetStoreScheduleByDateRangeItem{}, err
		}
		response = append(response, item)
	}
	return response, nil
}

type GetScheduleByIDItem struct {
	ID          int64       `db:"id"`
	WorkDate    pgtype.Date `db:"work_date"`
	Note        pgtype.Text `db:"note"`
	TimeSlotID  pgtype.Int8 `db:"time_slot_id"`
	StartTime   pgtype.Time `db:"start_time"`
	EndTime     pgtype.Time `db:"end_time"`
	IsAvailable pgtype.Bool `db:"is_available"`
}

func (r *ScheduleRepository) GetScheduleByID(ctx context.Context, scheduleID int64) ([]GetScheduleByIDItem, error) {
	query := `
		SELECT
			s.id,
			s.work_date,
			COALESCE(s.note, '') as note,
			ts.id as time_slot_id,
			ts.start_time,
			ts.end_time,
			ts.is_available
		FROM schedules s
		LEFT JOIN time_slots ts ON s.id = ts.schedule_id
		WHERE s.id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := []GetScheduleByIDItem{}
	for rows.Next() {
		var item GetScheduleByIDItem
		err := rows.Scan(
			&item.ID,
			&item.WorkDate,
			&item.Note,
			&item.TimeSlotID,
			&item.StartTime,
			&item.EndTime,
			&item.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		response = append(response, item)
	}

	return response, nil
}

func (r *ScheduleRepository) CheckScheduleExists(ctx context.Context, storeID int64, stylistID int64, workDate time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM schedules WHERE store_id = $1 AND stylist_id = $2 AND work_date = $3
		)
	`

	row := r.db.QueryRowxContext(ctx, query, storeID, stylistID, workDate)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

type UpdateScheduleParams struct {
	WorkDate *string
	Note     *string
}

type UpdateScheduleResponse struct {
	ID        int64
	StoreID   int64
	StylistID int64
	WorkDate  pgtype.Date
	Note      pgtype.Text
}

func (r *ScheduleRepository) UpdateSchedule(ctx context.Context, scheduleID int64, params UpdateScheduleParams) (UpdateScheduleResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{scheduleID}

	if params.WorkDate != nil {
		setParts = append(setParts, fmt.Sprintf("work_date = $%d", len(args)+1))
		args = append(args, *params.WorkDate)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	query := fmt.Sprintf(`
		UPDATE schedules
		SET %s
		WHERE id = $1
		RETURNING id, store_id, stylist_id, work_date, note
	`, strings.Join(setParts, ","))

	row := r.db.QueryRowxContext(ctx, query, args...)
	var response UpdateScheduleResponse
	err := row.Scan(
		&response.ID,
		&response.StoreID,
		&response.StylistID,
		&response.WorkDate,
		&response.Note,
	)
	if err != nil {
		return UpdateScheduleResponse{}, err
	}

	return response, nil
}
