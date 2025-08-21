package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllServicesByFilterParams struct {
	Name      *string
	IsAddon   *bool
	IsActive  *bool
	IsVisible *bool
	Limit     *int
	Offset    *int
	Sort      *[]string
}

type GetAllServicesByFilterItem struct {
	ID              int64              `db:"id"`
	Name            string             `db:"name"`
	Price           pgtype.Numeric     `db:"price"`
	DurationMinutes int32              `db:"duration_minutes"`
	IsAddon         pgtype.Bool        `db:"is_addon"`
	IsActive        pgtype.Bool        `db:"is_active"`
	IsVisible       pgtype.Bool        `db:"is_visible"`
	Note            pgtype.Text        `db:"note"`
	CreatedAt       pgtype.Timestamptz `db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `db:"updated_at"`
}

// GetAllServicesByFilter retrieves services with filtering, pagination and sorting
func (r *ServiceRepository) GetAllServicesByFilter(ctx context.Context, params GetAllServicesByFilterParams) (int, []GetAllServicesByFilterItem, error) {
	// where conditions
	whereConditions := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.IsAddon != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_addon = $%d", len(args)+1))
		args = append(args, *params.IsAddon)
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	if params.IsVisible != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_visible = $%d", len(args)+1))
		args = append(args, *params.IsVisible)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM services
		%s
	`, whereClause)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllServicesByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"created_at ASC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			price,
			duration_minutes,
			is_addon,
			is_active,
			is_visible,
			COALESCE(note, '') as note,
			created_at,
			updated_at
		FROM services
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllServicesByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateServiceParams struct {
	Name            *string
	Price           *int64
	DurationMinutes *int32
	IsAddon         *bool
	IsVisible       *bool
	IsActive        *bool
	Note            *string
}

type UpdateServiceResponse struct {
	ID              int64              `db:"id"`
	Name            string             `db:"name"`
	Price           pgtype.Numeric     `db:"price"`
	DurationMinutes int32              `db:"duration_minutes"`
	IsAddon         pgtype.Bool        `db:"is_addon"`
	IsVisible       pgtype.Bool        `db:"is_visible"`
	IsActive        pgtype.Bool        `db:"is_active"`
	Note            pgtype.Text        `db:"note"`
	CreatedAt       pgtype.Timestamptz `db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `db:"updated_at"`
}

func (r *ServiceRepository) UpdateService(ctx context.Context, serviceID int64, params UpdateServiceParams) (UpdateServiceResponse, error) {
	// set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", len(args)+1))
		args = append(args, *params.Price)
	}

	if params.DurationMinutes != nil {
		setParts = append(setParts, fmt.Sprintf("duration_minutes = $%d", len(args)+1))
		args = append(args, *params.DurationMinutes)
	}

	if params.IsAddon != nil {
		setParts = append(setParts, fmt.Sprintf("is_addon = $%d", len(args)+1))
		args = append(args, *params.IsAddon)
	}

	if params.IsVisible != nil {
		setParts = append(setParts, fmt.Sprintf("is_visible = $%d", len(args)+1))
		args = append(args, *params.IsVisible)
	}

	if params.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateServiceResponse{}, fmt.Errorf("no fields to update")
	}

	// Add WHERE clause
	args = append(args, serviceID)

	query := fmt.Sprintf(`
		UPDATE services
		SET %s
		WHERE id = $%d
		RETURNING id, name, price, duration_minutes, is_addon, is_visible, is_active, note, created_at, updated_at
	`, strings.Join(setParts, ", "), len(args))

	var result UpdateServiceResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateServiceResponse{}, fmt.Errorf("failed to update service: %w", err)
	}

	return result, nil
}
