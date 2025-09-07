package sqlx

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	SupplierID   pgtype.Int8        `db:"supplier_id"`
	SupplierName string             `db:"supplier_name"`
	PayerID      pgtype.Int8        `db:"payer_id"`
	PayerName    pgtype.Text        `db:"payer_name"`
	Category     pgtype.Text        `db:"category"`
	Amount       pgtype.Numeric     `db:"amount"`
	OtherFee     pgtype.Numeric     `db:"other_fee"`
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
			COALESCE(s.name, '') AS supplier_name,
			e.payer_id,
			COALESCE(su.username, '') AS payer_name,
			e.category,
			e.amount,
			e.other_fee,
			e.expense_date,
			e.note,
			e.is_reimbursed,
			e.reimbursed_at,
			e.created_at,
			e.updated_at
		FROM expenses e
		LEFT JOIN suppliers s ON e.supplier_id = s.id
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

// ------------------------------------------------------------------------------------------------

type UpdateStoreExpenseParams struct {
	SupplierID    *int64
	Category      *string
	Amount        *int64
	OtherFee      *int64
	ExpenseDate   *time.Time
	Note          *string
	PayerID       *int64
	PayerIDIsNone *bool
	IsReimbursed  *bool
	ReimbursedAt  *time.Time
}

type UpdateStoreExpenseResponse struct {
	ID int64 `db:"id"`
}

func (r *ExpenseRepository) UpdateStoreExpense(ctx context.Context, storeID, expenseID int64, req UpdateStoreExpenseParams) (UpdateStoreExpenseResponse, error) {
	// Set conditions
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if req.SupplierID != nil {
		setParts = append(setParts, fmt.Sprintf("supplier_id = $%d", len(args)+1))
		args = append(args, utils.Int64PtrToPgInt8(req.SupplierID))
	}

	if req.Category != nil && *req.Category != "" {
		setParts = append(setParts, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, utils.StringPtrToPgText(req.Category, false))
	}

	if req.Amount != nil {
		pgAmount, err := utils.Int64PtrToPgNumeric(req.Amount)
		if err != nil {
			return UpdateStoreExpenseResponse{}, fmt.Errorf("failed to convert amount: %w", err)
		}
		setParts = append(setParts, fmt.Sprintf("amount = $%d", len(args)+1))
		args = append(args, pgAmount)
	}

	if req.OtherFee != nil {
		pgOtherFee, err := utils.Int64PtrToPgNumeric(req.OtherFee)
		if err != nil {
			return UpdateStoreExpenseResponse{}, fmt.Errorf("failed to convert other fee: %w", err)
		}
		setParts = append(setParts, fmt.Sprintf("other_fee = $%d", len(args)+1))
		args = append(args, pgOtherFee)
	}

	if req.ExpenseDate != nil {
		pgDate := pgtype.Date{Time: *req.ExpenseDate, Valid: true}
		setParts = append(setParts, fmt.Sprintf("expense_date = $%d", len(args)+1))
		args = append(args, pgDate)
	}

	if req.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, utils.StringPtrToPgText(req.Note, true))
	}

	if req.PayerID != nil {
		setParts = append(setParts, fmt.Sprintf("payer_id = $%d", len(args)+1))
		args = append(args, utils.Int64PtrToPgInt8(req.PayerID))
	}

	if req.IsReimbursed != nil {
		setParts = append(setParts, fmt.Sprintf("is_reimbursed = $%d", len(args)+1))
		args = append(args, utils.BoolPtrToPgBool(req.IsReimbursed))
	}

	if req.ReimbursedAt != nil {
		setParts = append(setParts, fmt.Sprintf("reimbursed_at = $%d", len(args)+1))
		args = append(args, utils.TimePtrToPgTimestamptz(req.ReimbursedAt))
	}

	if req.PayerIDIsNone != nil && *req.PayerIDIsNone {
		// set payerId, isReimbursed, and reimbursedAt to nil
		setParts = append(setParts, fmt.Sprintf("payer_id = $%d", len(args)+1))
		args = append(args, nil)
		setParts = append(setParts, fmt.Sprintf("is_reimbursed = $%d", len(args)+1))
		args = append(args, nil)
		setParts = append(setParts, fmt.Sprintf("reimbursed_at = $%d", len(args)+1))
		args = append(args, nil)
	}

	// Check if there are any fields to update
	if len(setParts) == 1 {
		return UpdateStoreExpenseResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, expenseID, storeID)

	query := fmt.Sprintf(`
		UPDATE expenses
		SET %s
		WHERE id = $%d AND store_id = $%d
		RETURNING id
	`, strings.Join(setParts, ", "), len(args)-1, len(args))

	var result UpdateStoreExpenseResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateStoreExpenseResponse{}, fmt.Errorf("failed to execute update query: %w", err)
	}

	return result, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateStoreExpenseAmountTxParams struct {
	Amount int64
}

func (r *ExpenseRepository) UpdateStoreExpenseAmountTx(ctx context.Context, tx *sqlx.Tx, expenseID int64, params UpdateStoreExpenseAmountTxParams) error {
	query := `
		UPDATE expenses
		SET amount = $1
		WHERE id = $2
	`

	args := []interface{}{params.Amount, expenseID}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update expense amount: %w", err)
	}

	return nil
}
