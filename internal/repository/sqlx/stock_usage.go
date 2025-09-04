package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StockUsageRepository struct {
	db *sqlx.DB
}

func NewStockUsageRepository(db *sqlx.DB) *StockUsageRepository {
	return &StockUsageRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStockUsagesByFilterParams struct {
	ProductID *int64
	Name      *string
	IsInUse   *bool
	Limit     *int
	Offset    *int
	Sort      *[]string
}

type GetAllStockUsagesByFilterItem struct {
	ID           int64              `db:"id"`
	ProductID    int64              `db:"product_id"`
	ProductName  string             `db:"product_name"`
	Quantity     int32              `db:"quantity"`
	IsInUse      pgtype.Bool        `db:"is_in_use"`
	Expiration   pgtype.Date        `db:"expiration"`
	UsageStarted pgtype.Date        `db:"usage_started"`
	UsageEndedAt pgtype.Date        `db:"usage_ended_at"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}

// GetAllStockUsagesByFilter retrieves stock usages with filtering, pagination and sorting
func (r *StockUsageRepository) GetAllStockUsagesByFilter(ctx context.Context, storeID int64, params GetAllStockUsagesByFilterParams) (int, []GetAllStockUsagesByFilterItem, error) {
	// where conditions
	whereConditions := []string{"p.store_id = $1"}
	args := []interface{}{storeID}

	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.IsInUse != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("su.is_in_use = $%d", len(args)+1))
		args = append(args, *params.IsInUse)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM stock_usages su
		JOIN products p ON su.product_id = p.id
		%s
	`, whereClause)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllStockUsagesByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"su.created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "su.created_at",
		"updatedAt": "su.updated_at",
		"isInUse":   "su.is_in_use",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT
			su.id,
			p.id as product_id,
			p.name as product_name,
			su.quantity,
			su.is_in_use,
			su.expiration,
			su.usage_started,
			su.usage_ended_at,
			su.created_at,
			su.updated_at
		FROM stock_usages su
		JOIN products p ON su.product_id = p.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllStockUsagesByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute data query: %w", err)
	}

	return total, results, nil
}
