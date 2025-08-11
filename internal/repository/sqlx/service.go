package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ServiceRepositoryInterface interface {
	GetAllServiceByFilter(ctx context.Context, params GetAllServiceByFilterParams) (int, []GetAllServiceByFilterItem, error)
	UpdateService(ctx context.Context, serviceID int64, params UpdateServiceParams) (UpdateServiceResponse, error)
}

type ServiceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

// GetStoreServiceListModel represents the database model for admin service list queries
type GetAllServiceByFilterItem struct {
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

type GetAllServiceByFilterParams struct {
	Name      *string
	IsAddon   *bool
	IsActive  *bool
	IsVisible *bool
	Limit     *int
	Offset    *int
	Sort      *[]string
}

// GetStoreServiceList retrieves services for a specific store with admin filtering and pagination
func (r *ServiceRepository) GetAllServiceByFilter(ctx context.Context, params GetAllServiceByFilterParams) (int, []GetAllServiceByFilterItem, error) {
	// Set default values
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)

	// Set default sort values
	sort := utils.HandleSort([]string{"created_at", "updated_at", "is_active", "is_visible", "is_addon"}, "created_at", "ASC", params.Sort)

	// Build WHERE clause parts
	whereParts := []string{}
	args := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if params.Name != nil && *params.Name != "" {
		whereParts = append(whereParts, "name ILIKE :name")
		args["name"] = "%" + *params.Name + "%"
	}

	if params.IsAddon != nil {
		whereParts = append(whereParts, "is_addon = :is_addon")
		args["is_addon"] = *params.IsAddon
	}

	if params.IsActive != nil {
		whereParts = append(whereParts, "is_active = :is_active")
		args["is_active"] = *params.IsActive
	}

	if params.IsVisible != nil {
		whereParts = append(whereParts, "is_visible = :is_visible")
		args["is_visible"] = *params.IsVisible
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(whereParts, " AND "))
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM services
		%s
	`, whereClause)

	var total int
	rows, err := r.db.NamedQueryContext(ctx, countQuery, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, nil, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	// Main query with pagination
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
		LIMIT :limit OFFSET :offset
	`, whereClause, sort)

	var results []GetAllServiceByFilterItem
	rows, err = r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item GetAllServiceByFilterItem
		if err := rows.StructScan(&item); err != nil {
			return 0, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, item)
	}

	return total, results, nil
}

// ------------------------------------------------------------------------------------------------

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
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if params.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *params.Name)
		argIndex++
	}

	if params.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *params.Price)
		argIndex++
	}

	if params.DurationMinutes != nil {
		setParts = append(setParts, fmt.Sprintf("duration_minutes = $%d", argIndex))
		args = append(args, *params.DurationMinutes)
		argIndex++
	}

	if params.IsAddon != nil {
		setParts = append(setParts, fmt.Sprintf("is_addon = $%d", argIndex))
		args = append(args, *params.IsAddon)
		argIndex++
	}

	if params.IsVisible != nil {
		setParts = append(setParts, fmt.Sprintf("is_visible = $%d", argIndex))
		args = append(args, *params.IsVisible)
		argIndex++
	}

	if params.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", argIndex))
		args = append(args, *params.Note)
		argIndex++
	}

	// Add WHERE clause
	args = append(args, serviceID)
	whereClause := fmt.Sprintf("WHERE id = $%d", argIndex)

	query := fmt.Sprintf(`
		UPDATE services
		SET %s
		%s
		RETURNING id, name, price, duration_minutes, is_addon, is_visible, is_active, note, created_at, updated_at
	`, strings.Join(setParts, ", "), whereClause)

	var result UpdateServiceResponse
	err := r.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return UpdateServiceResponse{}, err
	}

	return result, nil
}
