package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TimeSlotTemplateRepositoryInterface defines the interface for time slot template repository
type TimeSlotTemplateRepositoryInterface interface {
	UpdateTimeSlotTemplate(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error)
	GetAllTimeSlotTemplateByFilter(ctx context.Context, params GetAllTimeSlotTemplateByFilterParams) (int, []GetAllTimeSlotTemplateByFilterItem, error)
}

type TimeSlotTemplateRepository struct {
	db *sqlx.DB
}

func NewTimeSlotTemplateRepository(db *sqlx.DB) *TimeSlotTemplateRepository {
	return &TimeSlotTemplateRepository{db: db}
}

// UpdateTimeSlotTemplate updates time slot template with dynamic fields
func (r *TimeSlotTemplateRepository) UpdateTimeSlotTemplate(ctx context.Context, templateID int64, req adminTimeSlotTemplateModel.UpdateTimeSlotTemplateRequest) (*adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := map[string]interface{}{
		"id": templateID,
	}

	if req.Name != nil {
		setParts = append(setParts, "name = :name")
		args["name"] = *req.Name
	}

	if req.Note != nil {
		setParts = append(setParts, "note = :note")
		args["note"] = *req.Note
	}

	query := fmt.Sprintf(`
		UPDATE time_slot_templates
		SET %s
		WHERE id = :id
		RETURNING id, name, note, updater, created_at, updated_at`,
		strings.Join(setParts, ", "))

	var result struct {
		ID        int64  `db:"id"`
		Name      string `db:"name"`
		Note      string `db:"note"`
		Updater   int64  `db:"updater"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	rows, err := r.db.NamedQuery(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update time slot template: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("time slot template not found")
	}

	if err := rows.StructScan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan result: %w", err)
	}

	return &adminTimeSlotTemplateModel.UpdateTimeSlotTemplateResponse{
		ID:   utils.FormatID(result.ID),
		Name: result.Name,
		Note: result.Note,
	}, nil
}

type GetAllTimeSlotTemplateByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Note      pgtype.Text        `db:"note"`
	Updater   int64              `db:"updater"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

type GetAllTimeSlotTemplateByFilterParams struct {
	Name   *string
	Limit  *int
	Offset *int
	Sort   *[]string
}

// GetTimeSlotTemplateList retrieves time slot templates with pagination and name filtering
func (r *TimeSlotTemplateRepository) GetAllTimeSlotTemplateByFilter(ctx context.Context, params GetAllTimeSlotTemplateByFilterParams) (int, []GetAllTimeSlotTemplateByFilterItem, error) {
	// Set pagination defaults
	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	sort := utils.HandleSort([]string{"created_at", "updated_at", "name"}, "updated_at", "DESC", params.Sort)

	// Build WHERE clause
	whereParts := []string{}
	args := []interface{}{}

	// Add name filter if provided
	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Query for total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM time_slot_templates %s`, whereClause)
	var total int
	row := r.db.QueryRowContext(ctx, countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return 0, nil, fmt.Errorf("count time slot templates failed: %w", err)
	}

	limitIdx := len(args) + 1
	offsetIdx := limitIdx + 1

	// Query for templates
	templatesQuery := fmt.Sprintf(`
		SELECT
			id,
			name,
			COALESCE(note, '') as note,
			updater,
			created_at,
			updated_at
		FROM time_slot_templates
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause,
		sort,
		limitIdx,
		offsetIdx,
	)

	argsWithLimit := append(args, limit, offset)

	var results []GetAllTimeSlotTemplateByFilterItem
	rows, err := r.db.QueryxContext(ctx, templatesQuery, argsWithLimit...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to query time slot templates: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result GetAllTimeSlotTemplateByFilterItem
		if err := rows.StructScan(&result); err != nil {
			return 0, nil, fmt.Errorf("failed to scan time slot template: %w", err)
		}
		results = append(results, result)
	}

	return total, results, nil
}
