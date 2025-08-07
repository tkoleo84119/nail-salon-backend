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
	GetStoreScheduleByDateRange(ctx context.Context, storeID int64, startDate time.Time, endDate time.Time, params GetStoreScheduleByDateRangeParams) ([]GetStoreScheduleByDateRangeItem, error)
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

// ------------------------------------------------------------------------------------------------

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
