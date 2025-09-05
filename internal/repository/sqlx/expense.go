package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ExpenseRepository struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllStoreExpensesByFilterParams struct {
	Category     *string
	SupplierID   *int64
	PayerID      *int64
	IsReimbursed *bool
	Limit        *int
	Offset       *int
	Sort         *[]string
}

type GetAllStoreExpensesByFilterItem struct {
	ID           int64              `db:"id"`
	SupplierID   int64              `db:"supplier_id"`
	SupplierName string             `db:"supplier_name"`
	PayerID      pgtype.Int8        `db:"payer_id"`
	PayerName    pgtype.Text        `db:"payer_name"`
	Category     pgtype.Text        `db:"category"`
	Amount       pgtype.Numeric     `db:"amount"`
	ExpenseDate  pgtype.Date        `db:"expense_date"`
	Note         pgtype.Text        `db:"note"`
	IsReimbursed pgtype.Bool        `db:"is_reimbursed"`
	ReimbursedAt pgtype.Timestamptz `db:"reimbursed_at"`
	CreatedAt    pgtype.Timestamptz `db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `db:"updated_at"`
}

func (r *ExpenseRepository) GetAllStoreExpensesByFilter(ctx context.Context, storeID int64, params GetAllStoreExpensesByFilterParams) (int, []GetAllStoreExpensesByFilterItem, error) {
	whereConditions := []string{"e.store_id = $1"}
	args := []interface{}{storeID}

	if params.Category != nil && *params.Category != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("e.category = $%d", len(args)+1))
		args = append(args, *params.Category)
	}

	if params.SupplierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.supplier_id = $%d", len(args)+1))
		args = append(args, *params.SupplierID)
	}

	if params.PayerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.payer_id = $%d", len(args)+1))
		args = append(args, *params.PayerID)
	}

	if params.IsReimbursed != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("e.is_reimbursed = $%d", len(args)+1))
		args = append(args, *params.IsReimbursed)
	}

	whereClause := "WHERE " + strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM expenses e
		%s
	`, whereClause)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}

	if total == 0 {
		return 0, []GetAllStoreExpensesByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"e.created_at DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"createdAt":    "e.created_at",
		"updatedAt":    "e.updated_at",
		"category":     "e.category",
		"supplierId":   "e.supplier_id",
		"payerId":      "e.payer_id",
		"isReimbursed": "e.is_reimbursed",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query with JOINs
	query := fmt.Sprintf(`
		SELECT
			e.id,
			e.supplier_id,
			s.name AS supplier_name,
			e.payer_id,
			su.name AS payer_name,
			e.category,
			e.amount,
			e.expense_date,
			e.note,
			e.is_reimbursed,
			e.reimbursed_at,
			e.created_at,
			e.updated_at
		FROM expenses e
		INNER JOIN suppliers s ON e.supplier_id = s.id
		LEFT JOIN staff_users su ON e.payer_id = su.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, limitIndex, offsetIndex)

	var results []GetAllStoreExpensesByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}
