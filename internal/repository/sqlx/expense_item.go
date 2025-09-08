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

type ExpenseItemRepository struct {
	db *sqlx.DB
}

func NewExpenseItemRepository(db *sqlx.DB) *ExpenseItemRepository {
	return &ExpenseItemRepository{
		db: db,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type UpdateStoreExpenseItemParams struct {
	ProductID       *int64
	Quantity        *int64
	Price           *int64
	ExpirationDate  *time.Time
	IsArrived       *bool
	ArrivalDate     *time.Time
	StorageLocation *string
	Note            *string
}

type UpdateStoreExpenseItemResponse struct {
	ID       int64          `db:"id"`
	Quantity int32          `db:"quantity"`
	Price    pgtype.Numeric `db:"price"`
}

func (r *ExpenseItemRepository) UpdateStoreExpenseItemTx(ctx context.Context, tx *sqlx.Tx, storeID, expenseID, expenseItemID int64, params UpdateStoreExpenseItemParams) (UpdateStoreExpenseItemResponse, error) {
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}

	if params.ProductID != nil {
		setParts = append(setParts, fmt.Sprintf("product_id = $%d", len(args)+1))
		args = append(args, *params.ProductID)
	}

	if params.Quantity != nil {
		setParts = append(setParts, fmt.Sprintf("quantity = $%d", len(args)+1))
		args = append(args, *params.Quantity)
	}

	if params.Price != nil {
		pgPrice, err := utils.Int64PtrToPgNumeric(params.Price)
		if err != nil {
			return UpdateStoreExpenseItemResponse{}, fmt.Errorf("failed to convert price: %w", err)
		}
		setParts = append(setParts, fmt.Sprintf("price = $%d", len(args)+1))
		args = append(args, pgPrice)
	}

	if params.ExpirationDate != nil {
		pgDate := utils.TimePtrToPgDate(params.ExpirationDate)
		setParts = append(setParts, fmt.Sprintf("expiration_date = $%d", len(args)+1))
		args = append(args, pgDate)
	}

	if params.IsArrived != nil {
		setParts = append(setParts, fmt.Sprintf("is_arrived = $%d", len(args)+1))
		args = append(args, utils.BoolPtrToPgBool(params.IsArrived))
	}

	if params.ArrivalDate != nil {
		pgDate := utils.TimePtrToPgDate(params.ArrivalDate)
		setParts = append(setParts, fmt.Sprintf("arrival_date = $%d", len(args)+1))
		args = append(args, pgDate)
	}

	if params.StorageLocation != nil {
		setParts = append(setParts, fmt.Sprintf("storage_location = $%d", len(args)+1))
		args = append(args, utils.StringPtrToPgText(params.StorageLocation, true))
	}

	if params.Note != nil {
		setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
		args = append(args, utils.StringPtrToPgText(params.Note, true))
	}

	if len(setParts) == 1 {
		return UpdateStoreExpenseItemResponse{}, fmt.Errorf("no fields to update")
	}

	args = append(args, expenseItemID)

	query := fmt.Sprintf(`
		UPDATE expense_items
		SET %s
		WHERE id = $%d
		RETURNING id, quantity, price
	`, strings.Join(setParts, ", "), len(args))

	var result UpdateStoreExpenseItemResponse
	if err := tx.GetContext(ctx, &result, query, args...); err != nil {
		return UpdateStoreExpenseItemResponse{}, fmt.Errorf("failed to execute update query: %w", err)
	}

	return result, nil
}
