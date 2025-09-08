package sqlx

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type AccountTransactionRepository struct {
	db *sqlx.DB
}

func NewAccountTransactionRepository(db *sqlx.DB) *AccountTransactionRepository {
	return &AccountTransactionRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type GetAllAccountTransactionsByFilterParams struct {
	Limit  *int
	Offset *int
	Sort   *[]string
}

type GetAllAccountTransactionsByFilterItem struct {
	ID              int64              `db:"id"`
	TransactionDate pgtype.Timestamptz `db:"transaction_date"`
	Type            string             `db:"type"`
	Amount          pgtype.Numeric     `db:"amount"`
	Balance         pgtype.Numeric     `db:"balance"`
	Note            pgtype.Text        `db:"note"`
	CreatedAt       pgtype.Timestamptz `db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `db:"updated_at"`
}

func (r *AccountTransactionRepository) GetAllAccountTransactionsByFilter(ctx context.Context, accountID int64, params GetAllAccountTransactionsByFilterParams) (int, []GetAllAccountTransactionsByFilterItem, error) {
	// where conditions
	args := []interface{}{accountID}

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM account_transactions
		WHERE account_id = $1
	`

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute count query: %w", err)
	}
	if total == 0 {
		return 0, []GetAllAccountTransactionsByFilterItem{}, nil
	}

	// Pagination + Sorting
	limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
	defaultSortArr := []string{"transaction_date DESC"}
	sort := utils.HandleSortByMap(map[string]string{
		"transactionDate": "transaction_date",
	}, defaultSortArr, params.Sort)

	args = append(args, limit, offset)
	limitIndex := len(args) - 1
	offsetIndex := len(args)

	// Data query
	query := fmt.Sprintf(`
		SELECT id, transaction_date, type, amount, balance, note, created_at, updated_at
		FROM account_transactions
		WHERE account_id = $1
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, sort, limitIndex, offsetIndex)

	var results []GetAllAccountTransactionsByFilterItem
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return 0, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return total, results, nil
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateAccountTransactionParams struct {
	Note *string
}

type UpdateAccountTransactionResponse struct {
	ID int64 `db:"id"`
}

func (r *AccountTransactionRepository) UpdateAccountTransaction(ctx context.Context, id int64, params UpdateAccountTransactionParams) (UpdateAccountTransactionResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	// Dynamic SET conditions
	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, *params.Note)
	}

	// Check if there are fields to update
	if len(setParts) == 1 {
		return UpdateAccountTransactionResponse{}, fmt.Errorf("no fields to update")
	}

	// Add WHERE condition ID
	args = append(args, id)
	whereIndex := len(args)

	query := fmt.Sprintf(`
		UPDATE account_transactions
		SET %s
		WHERE id = $%d
		RETURNING id
	`, strings.Join(setParts, ", "), whereIndex)

	var result UpdateAccountTransactionResponse
	if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateAccountTransactionResponse{}, fmt.Errorf("failed to update account transaction: %w", err)
	}

	return result, nil
}
