package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// StylistRepositoryInterface defines the interface for stylist repository
type StylistRepositoryInterface interface {
	GetStoreAllStylistByFilter(ctx context.Context, storeID int64, params GetStoreAllStylistByFilterParams) (int, []GetStoreAllStylistByFilterItem, error)
	UpdateStylist(ctx context.Context, staffUserID int64, params UpdateStylistParams) (UpdateStylistResponse, error)
}

type StylistRepository struct {
	db *sqlx.DB
}

func NewStylistRepository(db *sqlx.DB) *StylistRepository {
	return &StylistRepository{db: db}
}

type GetStoreAllStylistByFilterParams struct {
	Name        *string
	IsIntrovert *bool
	IsActive    *bool
	Limit       *int
	Offset      *int
	Sort        *[]string
}

type GetStoreAllStylistByFilterItem struct {
	ID           int64       `db:"id"`
	StaffUserID  int64       `db:"staff_user_id"`
	Name         pgtype.Text `db:"name"`
	GoodAtShapes []string    `db:"good_at_shapes"`
	GoodAtColors []string    `db:"good_at_colors"`
	GoodAtStyles []string    `db:"good_at_styles"`
	IsIntrovert  pgtype.Bool `db:"is_introvert"`
	IsActive     pgtype.Bool `db:"is_active"`
}

// GetStoreStylistList retrieves stylists for a specific store with dynamic filtering
func (r *StylistRepository) GetStoreAllStylistByFilter(ctx context.Context, storeID int64, params GetStoreAllStylistByFilterParams) (int, []GetStoreAllStylistByFilterItem, error) {
	// Set default values
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)

	// Set default sort values
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt":   "s.created_at",
		"updatedAt":   "s.updated_at",
		"isIntrovert": "s.is_introvert",
		"name":        "s.name",
		"isActive":    "sf.is_active",
	}, "s.created_at", "ASC", params.Sort)

	// Build WHERE conditions
	whereParts := []string{"sfsa.store_id = $1"}
	args := []interface{}{storeID}

	// Add name filter (case-insensitive partial match)
	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, fmt.Sprintf("s.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	// Add is_introvert filter
	if params.IsIntrovert != nil {
		whereParts = append(whereParts, fmt.Sprintf("s.is_introvert = $%d", len(args)+1))
		args = append(args, *params.IsIntrovert)
	}

	// Add is_active filter
	if params.IsActive != nil {
		whereParts = append(whereParts, fmt.Sprintf("sf.is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "WHERE " + strings.Join(whereParts, " AND ")
	}

	// Count total records with same filtering conditions
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT s.id)
		FROM stylists s
		INNER JOIN staff_users sf ON s.staff_user_id = sf.id
		INNER JOIN staff_user_store_access sfsa ON sf.id = sfsa.staff_user_id
		%s
	`, whereClause)

	var total int
	row := r.db.QueryRowContext(ctx, countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return 0, nil, fmt.Errorf("count stylists failed: %w", err)
	}

	limitIdx := len(args) + 1
	offsetIdx := limitIdx + 1

	// Query for stylists with filtering
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
	`, whereClause, sort, limitIdx, offsetIdx)

	argsWithLimit := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, argsWithLimit...)
	if err != nil {
		return 0, nil, fmt.Errorf("query stylists failed: %w", err)
	}
	defer rows.Close()

	m := pgtype.NewMap()
	var stylists []GetStoreAllStylistByFilterItem
	for rows.Next() {
		var stylist GetStoreAllStylistByFilterItem
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

// ------------------------------------------------------------------------------------------------

type UpdateStylistParams struct {
	Name         *string   `db:"name"`
	GoodAtShapes *[]string `db:"good_at_shapes"`
	GoodAtColors *[]string `db:"good_at_colors"`
	GoodAtStyles *[]string `db:"good_at_styles"`
	IsIntrovert  *bool     `db:"is_introvert"`
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
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{staffUserID}
	argIndex := 2

	if params.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *params.Name)
		argIndex++
	}

	if params.GoodAtShapes != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_shapes = $%d", argIndex))
		args = append(args, *params.GoodAtShapes)
		argIndex++
	}

	if params.GoodAtColors != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_colors = $%d", argIndex))
		args = append(args, *params.GoodAtColors)
		argIndex++
	}

	if params.GoodAtStyles != nil {
		setParts = append(setParts, fmt.Sprintf("good_at_styles = $%d", argIndex))
		args = append(args, *params.GoodAtStyles)
		argIndex++
	}

	if params.IsIntrovert != nil {
		setParts = append(setParts, fmt.Sprintf("is_introvert = $%d", argIndex))
		args = append(args, *params.IsIntrovert)
		argIndex++
	}

	query := fmt.Sprintf(`
		UPDATE stylists
		SET %s
		WHERE staff_user_id = $1
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
	`, strings.Join(setParts, ", "))

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
