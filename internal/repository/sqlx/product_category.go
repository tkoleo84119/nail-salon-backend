package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ProductCategoryRepository struct {
	db *sqlx.DB
}

func NewProductCategoryRepository(db *sqlx.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllProductCategoriesByFilterParams struct {
	Name     *string
	IsActive *bool
	Limit    *int
	Offset   *int
	Sort     *[]string
}

type GetAllProductCategoriesByFilterItem struct {
	ID        int64              `db:"id"`
	Name      string             `db:"name"`
	IsActive  pgtype.Bool        `db:"is_active"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func (r *ProductCategoryRepository) GetAllProductCategoriesByFilter(ctx context.Context, params GetAllProductCategoriesByFilterParams) (int, []GetAllProductCategoriesByFilterItem, error) {
	// where conditions
	whereConditions := []string{}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Name+"%")
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM product_categories
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}

	if total == 0 {
		return 0, []GetAllProductCategoriesByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"isActive":  "is_active",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"name":      "name",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT id, name, is_active, created_at, updated_at
		FROM product_categories
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllProductCategoriesByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateProductCategoryParams struct {
	Name     *string
	IsActive *bool
}

type UpdateProductCategoryResponse struct {
	ID int64 `db:"id"`
}

func (r *ProductCategoryRepository) UpdateProductCategory(ctx context.Context, id int64, params UpdateProductCategoryParams) (UpdateProductCategoryResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.Name != nil && *params.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *params.Name)
	}

	if params.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *params.IsActive)
	}

	if len(setParts) == 1 {
		return UpdateProductCategoryResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE product_categories
		SET %s
		WHERE id = $%d
		RETURNING id
	`, strings.Join(setParts, ", "), len(args))

	var response UpdateProductCategoryResponse
	err := r.db.GetContext(ctx, &response, query, args...)
	if err != nil {
		return UpdateProductCategoryResponse{}, err
	}

	return response, nil
}
