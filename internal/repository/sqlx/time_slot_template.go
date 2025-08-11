package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// TimeSlotTemplateRepositoryInterface defines the interface for time slot template repository
type TimeSlotTemplateRepositoryInterface interface {
	GetAllTimeSlotTemplateByFilter(ctx context.Context, params GetAllTimeSlotTemplateByFilterParams) (int, []GetAllTimeSlotTemplateByFilterItem, error)
	UpdateTimeSlotTemplate(ctx context.Context, templateID int64, params UpdateTimeSlotTemplateParams) (UpdateTimeSlotTemplateResponse, error)
}

type TimeSlotTemplateRepository struct {
	db *sqlx.DB
}

func NewTimeSlotTemplateRepository(db *sqlx.DB) *TimeSlotTemplateRepository {
	return &TimeSlotTemplateRepository{db: db}
}

type GetAllTimeSlotTemplateByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Note      pgtype.Text        `db:"note"`
	Updater   pgtype.Int8        `db:"updater"`
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
	// Set default values
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)

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

// ------------------------------------------------------------------------------------------------

type UpdateTimeSlotTemplateParams struct {
	Name *string
	Note *string
}

type UpdateTimeSlotTemplateResponse struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Note      pgtype.Text        `db:"note"`
	Updater   pgtype.Int8        `db:"updater"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

// UpdateTimeSlotTemplate updates time slot template with dynamic fields
func (r *TimeSlotTemplateRepository) UpdateTimeSlotTemplate(ctx context.Context, templateID int64, params UpdateTimeSlotTemplateParams) (UpdateTimeSlotTemplateResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{templateID}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	query := fmt.Sprintf(`
		UPDATE time_slot_templates
		SET %s
		WHERE id = $1
		RETURNING id, name, note, updater, created_at, updated_at`,
		strings.Join(setParts, ", "))

	row := r.db.QueryRowxContext(ctx, query, args...)
	var result UpdateTimeSlotTemplateResponse
	if err := row.StructScan(&result); err != nil {
		return UpdateTimeSlotTemplateResponse{}, fmt.Errorf("failed to scan time slot template: %w", err)
	}

	return result, nil
}
