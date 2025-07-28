package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// ScheduleRepositoryInterface defines the interface for schedule repository
type ScheduleRepositoryInterface interface {
	GetStoreScheduleList(ctx context.Context, storeID int64, params GetStoreScheduleListParams) (*adminScheduleModel.GetScheduleListResponse, error)
}

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

type GetStoreScheduleListParams struct {
	StylistID   *int64
	StartDate   time.Time
	EndDate     time.Time
	IsAvailable *bool
	Limit       int
	Offset      int
}

// Database models for individual queries
type ScheduleModel struct {
	ID        int64  `db:"id"`
	WorkDate  string `db:"work_date"`
	Note      string `db:"note"`
	StylistID int64  `db:"stylist_id"`
}

type StylistModel struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type TimeSlotModel struct {
	ID          int64  `db:"id"`
	ScheduleID  int64  `db:"schedule_id"`
	StartTime   string `db:"start_time"`
	EndTime     string `db:"end_time"`
	IsAvailable bool   `db:"is_available"`
}

// GetStoreScheduleList retrieves schedules for a specific store using step-by-step queries
func (r *ScheduleRepository) GetStoreScheduleList(ctx context.Context, storeID int64, params GetStoreScheduleListParams) (*adminScheduleModel.GetScheduleListResponse, error) {
	// Query schedules with basic filtering
	scheduleWhereParts := []string{
		"store_id = :store_id",
		"work_date BETWEEN :start_date AND :end_date",
	}
	scheduleArgs := map[string]interface{}{
		"store_id":   storeID,
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
	}

	// Add stylist filter if provided
	if params.StylistID != nil {
		scheduleWhereParts = append(scheduleWhereParts, "stylist_id = :stylist_id")
		scheduleArgs["stylist_id"] = *params.StylistID
	}

	scheduleWhereClause := strings.Join(scheduleWhereParts, " AND ")
	scheduleQuery := fmt.Sprintf(`
		SELECT id, work_date, COALESCE(note, '') as note, stylist_id
		FROM schedules
		WHERE %s
		ORDER BY work_date ASC, id ASC
	`, scheduleWhereClause)

	var schedules []ScheduleModel
	rows, err := r.db.NamedQueryContext(ctx, scheduleQuery, scheduleArgs)
	if err != nil {
		return nil, fmt.Errorf("query schedules failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule ScheduleModel
		if err := rows.StructScan(&schedule); err != nil {
			return nil, fmt.Errorf("scan schedule failed: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if len(schedules) == 0 {
		return &adminScheduleModel.GetScheduleListResponse{Items: []adminScheduleModel.GetScheduleListItem{}}, nil
	}

	// Get all unique stylist IDs and query stylist info
	stylistIDSet := make(map[int64]bool)
	scheduleIDs := make([]int64, 0, len(schedules))
	for _, schedule := range schedules {
		stylistIDSet[schedule.StylistID] = true
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}

	stylistIDList := make([]int64, 0, len(stylistIDSet))
	for stylistID := range stylistIDSet {
		stylistIDList = append(stylistIDList, stylistID)
	}

	// Query stylists
	stylistMap := make(map[int64]StylistModel)
	if len(stylistIDList) > 0 {
		stylistQuery, stylistArgs, err := sqlx.In("SELECT id, name FROM stylists WHERE id IN (?)", stylistIDList)
		if err != nil {
			return nil, fmt.Errorf("prepare stylist query failed: %w", err)
		}
		stylistQuery = r.db.Rebind(stylistQuery)

		stylistRows, err := r.db.QueryContext(ctx, stylistQuery, stylistArgs...)
		if err != nil {
			return nil, fmt.Errorf("query stylists failed: %w", err)
		}
		defer stylistRows.Close()

		for stylistRows.Next() {
			var stylist StylistModel
			if err := stylistRows.Scan(&stylist.ID, &stylist.Name); err != nil {
				return nil, fmt.Errorf("scan stylist failed: %w", err)
			}
			stylistMap[stylist.ID] = stylist
		}
	}

	// Query time slots for all schedules
	timeSlotWhereParts := []string{"schedule_id IN (?)"}
	timeSlotQueryArgs := []interface{}{scheduleIDs}

	// Add availability filter if specified
	if params.IsAvailable != nil {
		timeSlotWhereParts = append(timeSlotWhereParts, "is_available = ?")
		timeSlotQueryArgs = append(timeSlotQueryArgs, *params.IsAvailable)
	}

	timeSlotWhereClause := strings.Join(timeSlotWhereParts, " AND ")
	timeSlotQuery, timeSlotArgs, err := sqlx.In(fmt.Sprintf(`
		SELECT id, schedule_id, start_time::text as start_time, end_time::text as end_time, is_available
		FROM time_slots
		WHERE %s
		ORDER BY start_time ASC
	`, timeSlotWhereClause), timeSlotQueryArgs...)
	if err != nil {
		return nil, fmt.Errorf("prepare time slot query failed: %w", err)
	}
	timeSlotQuery = r.db.Rebind(timeSlotQuery)

	var timeSlots []TimeSlotModel
	timeSlotRows, err := r.db.QueryContext(ctx, timeSlotQuery, timeSlotArgs...)
	if err != nil {
		return nil, fmt.Errorf("query time slots failed: %w", err)
	}
	defer timeSlotRows.Close()

	for timeSlotRows.Next() {
		var timeSlot TimeSlotModel
		if err := timeSlotRows.Scan(&timeSlot.ID, &timeSlot.ScheduleID, &timeSlot.StartTime, &timeSlot.EndTime, &timeSlot.IsAvailable); err != nil {
			return nil, fmt.Errorf("scan time slot failed: %w", err)
		}
		timeSlots = append(timeSlots, timeSlot)
	}

	// Group time slots by schedule ID
	timeSlotMap := make(map[int64][]TimeSlotModel)
	for _, timeSlot := range timeSlots {
		timeSlotMap[timeSlot.ScheduleID] = append(timeSlotMap[timeSlot.ScheduleID], timeSlot)
	}

	// Build response
	items := make([]adminScheduleModel.GetScheduleListItem, 0, len(schedules))
	for _, schedule := range schedules {
		stylist := stylistMap[schedule.StylistID]
		scheduleTimeSlots := timeSlotMap[schedule.ID]

		timeSlotItems := make([]adminScheduleModel.GetScheduleListTimeSlotInfo, 0, len(scheduleTimeSlots))
		for _, ts := range scheduleTimeSlots {
			timeSlotItems = append(timeSlotItems, adminScheduleModel.GetScheduleListTimeSlotInfo{
				ID:          utils.FormatID(ts.ID),
				StartTime:   ts.StartTime,
				EndTime:     ts.EndTime,
				IsAvailable: ts.IsAvailable,
			})
		}

		items = append(items, adminScheduleModel.GetScheduleListItem{
			ID:       utils.FormatID(schedule.ID),
			WorkDate: schedule.WorkDate,
			Stylist: adminScheduleModel.GetScheduleListStylistInfo{
				ID:   utils.FormatID(stylist.ID),
				Name: stylist.Name,
			},
			Note:      schedule.Note,
			TimeSlots: timeSlotItems,
		})
	}

	return &adminScheduleModel.GetScheduleListResponse{Items: items}, nil
}
