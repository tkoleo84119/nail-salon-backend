package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStoreProductsByFilterParams struct {
	BrandID             *int64
	CategoryID          *int64
	Name                *string
	LessThanSafetyStock *bool
	Limit               *int
	Offset              *int
	Sort                *[]string
}

type GetAllStoreProductsByFilterItem struct {
	ID              int64              `db:"id"`
	Name            string             `db:"name"`
	BrandID         int64              `db:"brand_id"`
	BrandName       string             `db:"brand_name"`
	CategoryID      int64              `db:"category_id"`
	CategoryName    string             `db:"category_name"`
	CurrentStock    int32              `db:"current_stock"`
	SafetyStock     pgtype.Int4        `db:"safety_stock"`
	Unit            pgtype.Text        `db:"unit"`
	StorageLocation pgtype.Text        `db:"storage_location"`
	Note            pgtype.Text        `db:"note"`
	CreatedAt       pgtype.Timestamptz `db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `db:"updated_at"`
}

func (r *ProductRepository) GetAllStoreProductsByFilter(ctx context.Context, storeID int64, params GetAllStoreProductsByFilterParams) (int, []GetAllStoreProductsByFilterItem, error) {
	whereConditions := []string{"p.store_id = $1"}
	args := []interface{}{storeID}

	if params.BrandID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.brand_id = $%d", len(args)+1))
		args = append(args, *params.BrandID)
	}

	if params.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.category_id = $%d", len(args)+1))
		args = append(args, *params.CategoryID)
	}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.LessThanSafetyStock != nil && *params.LessThanSafetyStock {
		whereConditions = append(whereConditions, "p.current_stock <= p.safety_stock")
	}

	whereClause := "WHERE " + strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM products p
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllStoreProductsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"p.created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt":  "p.created_at",
		"updatedAt":  "p.updated_at",
		"brandId":    "p.brand_id",
		"categoryId": "p.category_id",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query with JOINs
	query := fmt.Sprintf(`
		SELECT
			p.id,
			p.name,
			p.brand_id,
			b.name AS brand_name,
			p.category_id,
			pc.name AS category_name,
			p.current_stock,
			p.safety_stock,
			p.unit,
			p.storage_location,
			p.note,
			p.created_at,
			p.updated_at
		FROM products p
		INNER JOIN brands b ON p.brand_id = b.id
		INNER JOIN product_categories pc ON p.category_id = pc.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllStoreProductsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}
