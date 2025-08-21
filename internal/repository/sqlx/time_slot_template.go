package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type TimeSlotTemplateRepository struct {
	db *sqlx.DB
}

func NewTimeSlotTemplateRepository(db *sqlx.DB) *TimeSlotTemplateRepository {
	return &TimeSlotTemplateRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllTimeSlotTemplatesByFilterParams struct {
	Name   *string
	Limit  *int
	Offset *int
	Sort   *[]string
}

type GetAllTimeSlotTemplatesByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	Note      pgtype.Text        `db:"note"`
	Updater   pgtype.Int8        `db:"updater"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

// GetAllTimeSlotTemplatesByFilter retrieves time slot templates with pagination and name filtering
func (r *TimeSlotTemplateRepository) GetAllTimeSlotTemplatesByFilter(ctx context.Context, params GetAllTimeSlotTemplatesByFilterParams) (int, []GetAllTimeSlotTemplatesByFilterItem, error) {
	// WHERE clause
	whereParts := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM time_slot_templates %s`, whereClause)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("count time slot templates failed: %w", err)
	}
	if total == 0 {
		return 0, []GetAllTimeSlotTemplatesByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"updated_at ASC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"name":      "name",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
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
		limitIndex,
		offsetIndex,
	)

	var results []GetAllTimeSlotTemplatesByFilterItem
	if err := r.db.SelectContext(ctx, &results, templatesQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("query time slot templates failed: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

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
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateTimeSlotTemplateResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, templateID)

	query := fmt.Sprintf(`
		UPDATE time_slot_templates
		SET %s
		WHERE id = $%d
		RETURNING id, name, note, updater, created_at, updated_at
	`, strings.Join(setParts, ", "), len(args))

	var result UpdateTimeSlotTemplateResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateTimeSlotTemplateResponse{}, fmt.Errorf("update time slot template failed: %w", err)
	}

	return result, nil
}
