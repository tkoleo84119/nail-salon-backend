package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type BookingProductRepository struct {
	db *sqlx.DB
}

func NewBookingProductRepository(db *sqlx.DB) *BookingProductRepository {
	return &BookingProductRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllBookingProductsByFilterParams struct {
	BookingID int64
	Limit     *int
	Offset    *int
	Sort      *[]string
}

type GetAllBookingProductsByFilterItem struct {
	ProductID    int64              `db:"product_id"`
	ProductName  string             `db:"product_name"`
	BrandID      int64              `db:"brand_id"`
	BrandName    string             `db:"brand_name"`
	CategoryID   int64              `db:"category_id"`
	CategoryName string             `db:"category_name"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
}

func (r *BookingProductRepository) GetAllBookingProductsByFilter(ctx context.Context, params GetAllBookingProductsByFilterParams) (int, []GetAllBookingProductsByFilterItem, error) {
	// Build WHERE clause
	whereConditions := []string{"bp.booking_id = $1"}
	args := []interface{}{params.BookingID}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM booking_products bp
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to count booking products: %w", err)
	}
	if total == 0 {
		return 0, []GetAllBookingProductsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"bp.created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt": "bp.created_at",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query with JOINs
	query := fmt.Sprintf(`
		SELECT
			bp.product_id,
			p.name AS product_name,
			b.id AS brand_id,
			b.name AS brand_name,
			pc.id AS category_id,
			pc.name AS category_name,
			bp.created_at
		FROM booking_products bp
		INNER JOIN products p ON bp.product_id = p.id
		INNER JOIN brands b ON p.brand_id = b.id
		INNER JOIN product_categories pc ON p.category_id = pc.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllBookingProductsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}
