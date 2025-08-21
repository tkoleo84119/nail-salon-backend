package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StylistRepository struct {
	db *sqlx.DB
}

func NewStylistRepository(db *sqlx.DB) *StylistRepository {
	return &StylistRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStoreStylistsByFilterParams struct {
	Name        *string
	IsIntrovert *bool
	IsActive    *bool
	Limit       *int
	Offset      *int
	Sort        *[]string
}

type GetAllStoreStylistsByFilterItem struct {
	ID           int64       `db:"id"`
	StaffUserID  int64       `db:"staff_user_id"`
	Name         pgtype.Text `db:"name"`
	GoodAtShapes []string    `db:"good_at_shapes"`
	GoodAtColors []string    `db:"good_at_colors"`
	GoodAtStyles []string    `db:"good_at_styles"`
	IsIntrovert  pgtype.Bool `db:"is_introvert"`
	IsActive     pgtype.Bool `db:"is_active"`
}

// GetAllStoreStylistsByFilter retrieves stylists for a specific store with dynamic filtering
func (r *StylistRepository) GetAllStoreStylistsByFilter(ctx context.Context, storeID int64, params GetAllStoreStylistsByFilterParams) (int, []GetAllStoreStylistsByFilterItem, error) {
	// where conditions
	whereParts := []string{"sfsa.store_id = $1"}
	args := []interface{}{storeID}

	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, fmt.Sprintf("s.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.IsIntrovert != nil {
		whereParts = append(whereParts, fmt.Sprintf("s.is_introvert = $%d", len(args)+1))
		args = append(args, *params.IsIntrovert)
	}

	if params.IsActive != nil {
		whereParts = append(whereParts, fmt.Sprintf("sf.is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT s.id)
		FROM stylists s
		INNER JOIN staff_users sf ON s.staff_user_id = sf.id
		INNER JOIN staff_user_store_access sfsa ON sf.id = sfsa.staff_user_id
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("count stylists failed: %w", err)
	}
	if total == 0 {
		return 0, []GetAllStoreStylistsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"s.created_at ASC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt":   "s.created_at",
		"updatedAt":   "s.updated_at",
		"isIntrovert": "s.is_introvert",
		"name":        "s.name",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT
			s.id,
			s.staff_user_id,
			s.name,
			COALESCE(s.good_at_shapes, '{}'::text[]) AS good_at_shapes,
			COALESCE(s.good_at_colors, '{}'::text[]) AS good_at_colors,
			COALESCE(s.good_at_styles, '{}'::text[]) AS good_at_styles,
			s.is_introvert,
			sf.is_active AS is_active
		FROM stylists s
		INNER JOIN staff_users sf ON s.staff_user_id = sf.id
		INNER JOIN staff_user_store_access sfsa ON sf.id = sfsa.staff_user_id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("query stylists failed: %w", err)
	}
	defer rows.Close()

	m := pgtype.NewMap()
	var stylists []GetAllStoreStylistsByFilterItem
	for rows.Next() {
		var stylist GetAllStoreStylistsByFilterItem
		if err := rows.Scan(
			&stylist.ID,
			&stylist.StaffUserID,
			&stylist.Name,
			m.SQLScanner(&stylist.GoodAtShapes),
			m.SQLScanner(&stylist.GoodAtColors),
			m.SQLScanner(&stylist.GoodAtStyles),
			&stylist.IsIntrovert,
			&stylist.IsActive,
		); err != nil {
			return 0, nil, fmt.Errorf("scan stylist failed: %w", err)
		}
		stylists = append(stylists, stylist)
	}

	return total, stylists, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateStylistParams struct {
	Name         *string
	GoodAtShapes *[]string
	GoodAtColors *[]string
	GoodAtStyles *[]string
	IsIntrovert  *bool
}

type UpdateStylistResponse struct {
	ID           int64              `db:"id"`
	StaffUserID  int64              `db:"staff_user_id"`
	Name         pgtype.Text        `db:"name"`
	GoodAtShapes []string           `db:"good_at_shapes"`
	GoodAtColors []string           `db:"good_at_colors"`
	GoodAtStyles []string           `db:"good_at_styles"`
	IsIntrovert  pgtype.Bool        `db:"is_introvert"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}

// UpdateStylist updates stylist with dynamic fields
func (r *StylistRepository) UpdateStylist(ctx context.Context, staffUserID int64, params UpdateStylistParams) (UpdateStylistResponse, error) {
	// Set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.GoodAtShapes != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_shapes = $%d", len(args)+1))
		args = append(args, *params.GoodAtShapes)
	}

	if params.GoodAtColors != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_colors = $%d", len(args)+1))
		args = append(args, *params.GoodAtColors)
	}

	if params.GoodAtStyles != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_styles = $%d", len(args)+1))
		args = append(args, *params.GoodAtStyles)
	}

	if params.IsIntrovert != nil {
		setParts = append(setParts, fmt.Sprintf("is_introvert = $%d", len(args)+1))
		args = append(args, *params.IsIntrovert)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateStylistResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, staffUserID)

	// Data query
	query := fmt.Sprintf(`
		UPDATE stylists
		SET %s
		WHERE staff_user_id = $%d
		RETURNING
			id,
			staff_user_id,
			name,
			COALESCE(good_at_shapes, '{}'::text[]) AS good_at_shapes,
			COALESCE(good_at_colors, '{}'::text[]) AS good_at_colors,
			COALESCE(good_at_styles, '{}'::text[]) AS good_at_styles,
			is_introvert,
			created_at,
			updated_at
	`, strings.Join(setParts, ", "), len(args))

	row := r.db.QueryRowxContext(ctx, query, args...)
	m := pgtype.NewMap()

	var result UpdateStylistResponse

	err := row.Scan(
		&result.ID,
		&result.StaffUserID,
		&result.Name,
		m.SQLScanner(&result.GoodAtShapes),
		m.SQLScanner(&result.GoodAtColors),
		m.SQLScanner(&result.GoodAtStyles),
		&result.IsIntrovert,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return UpdateStylistResponse{}, fmt.Errorf("scan result failed: %w", err)
	}

	return result, nil
}
