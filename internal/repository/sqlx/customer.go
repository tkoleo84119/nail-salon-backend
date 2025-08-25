package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CustomerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllCustomersByFilterParams struct {
	Name          *string
	LineName      *string
	Phone         *string
	Level         *string
	IsBlacklisted *bool
	MinPastDays   *int
	Limit         *int
	Offset        *int
	Sort          *[]string
}

type GetAllCustomersByFilterItem struct {
	ID            int64              `db:"id"`
	Name          string             `db:"name"`
	LineName      pgtype.Text        `db:"line_name"`
	Phone         string             `db:"phone"`
	Birthday      pgtype.Date        `db:"birthday"`
	City          pgtype.Text        `db:"city"`
	Level         pgtype.Text        `db:"level"`
	IsBlacklisted pgtype.Bool        `db:"is_blacklisted"`
	LastVisitAt   pgtype.Timestamptz `db:"last_visit_at"`
	UpdatedAt     pgtype.Timestamptz `db:"updated_at"`
}

// GetAllCustomersByFilter retrieves all customers with filtering, pagination and sorting
func (r *CustomerRepository) GetAllCustomersByFilter(ctx context.Context, params GetAllCustomersByFilterParams) (int, []GetAllCustomersByFilterItem, error) {
	// where conditions
	whereConditions := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.LineName != nil && *params.LineName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("line_name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.LineName+"%")
	}

	if params.Phone != nil && *params.Phone != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("phone ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Phone+"%")
	}

	if params.Level != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("level = $%d", len(args)+1))
		args = append(args, *params.Level)
	}

	if params.IsBlacklisted != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_blacklisted = $%d", len(args)+1))
		args = append(args, *params.IsBlacklisted)
	}

	if params.MinPastDays != nil && *params.MinPastDays > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("(last_visit_at IS NOT NULL AND last_visit_at < NOW() - ($%d * INTERVAL '1 day'))", len(args)+1))
		args = append(args, *params.MinPastDays)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM customers
		%s`, whereClause)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllCustomersByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"last_visit_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"level":         "level",
		"isBlacklisted": "is_blacklisted",
		"lastVisitAt":   "last_visit_at",
		"createdAt":     "created_at",
		"updatedAt":     "updated_at",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	dataQuery := fmt.Sprintf(`
		SELECT
			id, name, line_name, phone, birthday, city,
			level, is_blacklisted, last_visit_at, updated_at
		FROM customers
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllCustomersByFilterItem
	if err := r.db.SelectContext(ctx, &results, dataQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateCustomerParams struct {
	Name           *string
	LineName       *string
	Phone          *string
	Birthday       *string
	Email          *string
	City           *string
	FavoriteShapes *[]string
	FavoriteColors *[]string
	FavoriteStyles *[]string
	IsIntrovert    *bool
	CustomerNote   *string
	StoreNote      *string
	Level          *string
	IsBlacklisted  *bool
}

type UpdateCustomerResponse struct {
	ID             int64              `db:"id"`
	Name           string             `db:"name"`
	LineName       pgtype.Text        `db:"line_name"`
	Phone          string             `db:"phone"`
	Birthday       pgtype.Date        `db:"birthday"`
	Email          pgtype.Text        `db:"email"`
	City           pgtype.Text        `db:"city"`
	FavoriteShapes []string           `db:"favorite_shapes"`
	FavoriteColors []string           `db:"favorite_colors"`
	FavoriteStyles []string           `db:"favorite_styles"`
	IsIntrovert    pgtype.Bool        `db:"is_introvert"`
	ReferralSource []string           `db:"referral_source"`
	Referrer       pgtype.Text        `db:"referrer"`
	CustomerNote   pgtype.Text        `db:"customer_note"`
	StoreNote      pgtype.Text        `db:"store_note"`
	Level          pgtype.Text        `db:"level"`
	IsBlacklisted  pgtype.Bool        `db:"is_blacklisted"`
	LastVisitAt    pgtype.Timestamptz `db:"last_visit_at"`
	CreatedAt      pgtype.Timestamptz `db:"created_at"`
	UpdatedAt      pgtype.Timestamptz `db:"updated_at"`
}

func (r *CustomerRepository) UpdateCustomer(ctx context.Context, customerID int64, params UpdateCustomerParams) (UpdateCustomerResponse, error) {
	// set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.LineName != nil && *params.LineName != "" {
		setParts = append(setParts, fmt.Sprintf("line_name = $%d", len(args)+1))
		args = append(args, *params.LineName)
	}

	if params.Phone != nil && *params.Phone != "" {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", len(args)+1))
		args = append(args, *params.Phone)
	}

	if params.Birthday != nil {
		setParts = append(setParts, fmt.Sprintf("birthday = $%d", len(args)+1))
		args = append(args, *params.Birthday)
	}

	if params.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, *params.Email)
	}

	if params.City != nil {
		setParts = append(setParts, fmt.Sprintf("city = $%d", len(args)+1))
		args = append(args, *params.City)
	}

	if params.FavoriteShapes != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_shapes = $%d", len(args)+1))
		args = append(args, *params.FavoriteShapes)
	}

	if params.FavoriteColors != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_colors = $%d", len(args)+1))
		args = append(args, *params.FavoriteColors)
	}

	if params.FavoriteStyles != nil {
		setParts = append(setParts, fmt.Sprintf("favorite_styles = $%d", len(args)+1))
		args = append(args, *params.FavoriteStyles)
	}

	if params.IsIntrovert != nil {
		setParts = append(setParts, fmt.Sprintf("is_introvert = $%d", len(args)+1))
		args = append(args, *params.IsIntrovert)
	}

	if params.StoreNote != nil {
		setParts = append(setParts, fmt.Sprintf("store_note = $%d", len(args)+1))
		args = append(args, *params.StoreNote)
	}

	if params.CustomerNote != nil {
		setParts = append(setParts, fmt.Sprintf("customer_note = $%d", len(args)+1))
		args = append(args, *params.CustomerNote)
	}

	if params.Level != nil {
		setParts = append(setParts, fmt.Sprintf("level = $%d", len(args)+1))
		args = append(args, *params.Level)
	}

	if params.IsBlacklisted != nil {
		setParts = append(setParts, fmt.Sprintf("is_blacklisted = $%d", len(args)+1))
		args = append(args, *params.IsBlacklisted)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateCustomerResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, customerID)

	// Data query
	query := fmt.Sprintf(`
		UPDATE customers
		SET %s
		WHERE id = $%d
		RETURNING
			id,
			name,
			line_name,
			phone,
			birthday,
			email,
			city,
			COALESCE(favorite_shapes, '{}'::text[]) AS favorite_shapes,
			COALESCE(favorite_colors, '{}'::text[]) AS favorite_colors,
			COALESCE(favorite_styles, '{}'::text[]) AS favorite_styles,
			is_introvert,
			COALESCE(referral_source, '{}'::text[]) AS referral_source,
			referrer,
			customer_note,
			store_note,
			level,
			is_blacklisted,
			last_visit_at,
			created_at,
			updated_at
		`,
		strings.Join(setParts, ", "), len(args))

	row := r.db.QueryRowxContext(ctx, query, args...)
	m := pgtype.NewMap()

	var result UpdateCustomerResponse
	err := row.Scan(
		&result.ID,
		&result.Name,
		&result.LineName,
		&result.Phone,
		&result.Birthday,
		&result.Email,
		&result.City,
		m.SQLScanner(&result.FavoriteShapes),
		m.SQLScanner(&result.FavoriteColors),
		m.SQLScanner(&result.FavoriteStyles),
		&result.IsIntrovert,
		m.SQLScanner(&result.ReferralSource),
		&result.Referrer,
		&result.CustomerNote,
		&result.StoreNote,
		&result.Level,
		&result.IsBlacklisted,
		&result.LastVisitAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return UpdateCustomerResponse{}, fmt.Errorf("scan result failed: %w", err)
	}

	return result, nil
}
