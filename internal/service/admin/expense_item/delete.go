package adminExpenseItem

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
}

func NewDelete(queries *dbgen.Queries, db *pgxpool.Pool) DeleteInterface {
	return &Delete{
		queries: queries,
		db:      db,
	}
}

func (s *Delete) Delete(ctx context.Context, storeID, expenseID, expenseItemID int64, role string, creatorID int64,creatorStoreIDs []int64) (*adminExpenseItemModel.DeleteResponse, error) {
	// Check store access permission
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// Check if expense exists
	expense, err := s.queries.GetStoreExpenseByID(ctx, dbgen.GetStoreExpenseByIDParams{
		ID:      expenseID,
		StoreID: storeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense", err)
	}

	// Check if expense is reimbursed
	if expense.IsReimbursed.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseReimbursedNotAllowToDeleteItem)
	}

	// Check if expense item exists
	expenseItem, err := s.queries.GetStoreExpenseItemByID(ctx, dbgen.GetStoreExpenseItemByIDParams{
		ID:        expenseItemID,
		ExpenseID: expenseID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense item", err)
	}

	// Check if expense item is arrived
	if expenseItem.IsArrived.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemArrivedNotAllowToDelete)
	}

	// calculate new total amount
	oldAmount, err := utils.PgNumericToFloat64(expense.Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert old amount to float64", err)
	}
	itemPrice, err := utils.PgNumericToFloat64(expenseItem.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert item price to float64", err)
	}
	newAmount := oldAmount - (itemPrice * float64(expenseItem.Quantity))
	newAmountNumeric, err := utils.Float64PtrToPgNumeric(&newAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert new amount to pgnumeric", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Delete expense item
	err = qtx.DeleteStoreExpenseItem(ctx, dbgen.DeleteStoreExpenseItemParams{
		ID:        expenseItemID,
		ExpenseID: expenseID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete expense item", err)
	}

	// Update expense amount
	err = qtx.UpdateStoreExpenseAmount(ctx, dbgen.UpdateStoreExpenseAmountParams{
		Amount:  newAmountNumeric,
		Updater: utils.Int64PtrToPgInt8(&creatorID),
		ID:      expenseID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense amount", err)
	}

	return &adminExpenseItemModel.DeleteResponse{
		Deleted: utils.FormatID(expenseItemID),
	}, nil
}
