package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetStoreSchedulesByDateRangeParams struct {
	StylistID   *[]int64
	IsAvailable *bool
}

type GetStoreSchedulesByDateRangeItem struct {
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

// GetStoreSchedulesByDateRange retrieves schedules for a specific store using step-by-step queries
func (r *ScheduleRepository) GetStoreSchedulesByDateRange(ctx context.Context, storeID int64, startDate time.Time, endDate time.Time, params GetStoreSchedulesByDateRangeParams) ([]GetStoreSchedulesByDateRangeItem, error) {
	whereParts := []string{
		"s.store_id = $1",
		"s.work_date BETWEEN $2 AND $3",
	}
	args := []interface{}{storeID, startDate, endDate}

	// Add stylist filter if provided
	if params.StylistID != nil && len(*params.StylistID) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("s.stylist_id = ANY($%d)", len(args)+1))
		args = append(args, *params.StylistID)
	}

	if params.IsAvailable != nil {
		whereParts = append(whereParts, fmt.Sprintf("ts.is_available = $%d", len(args)+1))
		args = append(args, *params.IsAvailable)
	}

	whereClause := strings.Join(whereParts, " AND ")

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

	var results []GetStoreSchedulesByDateRangeItem
	err := r.db.SelectContext(ctx, &results, scheduleQuery, args...)
	if err != nil {
		return []GetStoreSchedulesByDateRangeItem{}, fmt.Errorf("failed to execute data query: %w", err)
	}

	return results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateScheduleParams struct {
	WorkDate pgtype.Date
	Note     *string
}

type UpdateScheduleResponse struct {
	ID        int64       `db:"id"`
	StoreID   int64       `db:"store_id"`
	StylistID int64       `db:"stylist_id"`
	WorkDate  pgtype.Date `db:"work_date"`
	Note      pgtype.Text `db:"note"`
}

func (r *ScheduleRepository) UpdateSchedule(ctx context.Context, scheduleID int64, params UpdateScheduleParams) (UpdateScheduleResponse, error) {
	// set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.WorkDate.Valid {
		setParts = append(setParts, fmt.Sprintf("work_date = $%d", len(args)+1))
		args = append(args, params.WorkDate)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateScheduleResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, scheduleID)

	query := fmt.Sprintf(`
		UPDATE schedules
		SET %s
		WHERE id = $%d
		RETURNING id, store_id, stylist_id, work_date, note
	`, strings.Join(setParts, ","), len(args))

	var response UpdateScheduleResponse
	if err := r.db.GetContext(ctx, &response, query, args...); err != nil {
		return UpdateScheduleResponse{}, fmt.Errorf("failed to update schedule: %w", err)
	}

	return response, nil
}
