package adminExpenseItem

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
	db      *sqlx.DB
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, db *sqlx.DB) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
		db:      db,
	}
}

func (s *Update) Update(ctx context.Context, storeID, expenseID, expenseItemID int64, req adminExpenseItemModel.UpdateParsedRequest, updaterID int64, role string, creatorStoreIDs []int64) (*adminExpenseItemModel.UpdateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// when pass isArrived is true, but not pass arrivalDate
	if req.IsArrived != nil && *req.IsArrived && req.ArrivalDate == nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemNotAllowPassIsArrivedTrueWithoutArrivalDate)
	}

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

	// if expense is reimbursed, not allow to update product id, quantity, price
	if expense.IsReimbursed.Valid && expense.IsReimbursed.Bool {
		if req.ProductID != nil || req.Quantity != nil || req.Price != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseReimbursedNotAllowToUpdateProductInfo)
		}
	}

	oldExpenseItem, err := s.queries.GetStoreExpenseItemByID(ctx, dbgen.GetStoreExpenseItemByIDParams{
		ID:        expenseItemID,
		ExpenseID: expenseID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense item", err)
	}

	if err := s.validateArrivedItemRestrictions(oldExpenseItem, req); err != nil {
		return nil, err
	}

	if req.ProductID != nil {
		product, err := s.queries.GetProductByID(ctx, *req.ProductID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
			}
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get product", err)
		}
		if product.StoreID != storeID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
		}
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	updatedExpenseItem, err := s.repo.ExpenseItem.UpdateStoreExpenseItemTx(ctx, tx, storeID, expenseID, expenseItemID, sqlxRepo.UpdateStoreExpenseItemParams{
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		Price:           req.Price,
		ExpirationDate:  req.ExpirationDate,
		IsArrived:       req.IsArrived,
		ArrivalDate:     req.ArrivalDate,
		StorageLocation: req.StorageLocation,
		Note:            req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense item", err)
	}

	if req.Price != nil || req.Quantity != nil {
		if err := s.updateExpenseAmount(ctx, tx, expense, oldExpenseItem, updatedExpenseItem, expenseID, updaterID); err != nil {
			return nil, err
		}
	} else {
		err = s.repo.Expense.UpdateStoreExpenseUpdaterTx(ctx, tx, expenseID, sqlxRepo.UpdateStoreExpenseUpdaterTxParams{
			Updater: updaterID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense updater", err)
		}
	}

	// update product stock
	if req.IsArrived != nil && *req.IsArrived && !oldExpenseItem.IsArrived.Bool {
		if err := s.updateProductStockForArrival(ctx, tx, oldExpenseItem, req); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminExpenseItemModel.UpdateResponse{
		ID: utils.FormatID(updatedExpenseItem.ID),
	}, nil
}

func (s *Update) validateArrivedItemRestrictions(oldExpenseItem dbgen.GetStoreExpenseItemByIDRow, req adminExpenseItemModel.UpdateParsedRequest) error {
	// if product is not arrived, no restriction
	if !oldExpenseItem.IsArrived.Bool {
		return nil
	}

	// 1. not allow to change back to not arrived
	if req.IsArrived != nil && !*req.IsArrived {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemArrivedNotAllowToTurnBack)
	}

	// 2. not allow to change product id
	if req.ProductID != nil && *req.ProductID != oldExpenseItem.ProductID {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemArrivedNotAllowToChangeProductID)
	}

	// 3. not allow to change quantity
	if req.Quantity != nil && *req.Quantity != int64(oldExpenseItem.Quantity) {
		return errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemArrivedNotAllowToChangeQuantity)
	}

	// 4. not allow to change price
	if req.Price != nil {
		oldPrice, err := utils.PgNumericToFloat64(oldExpenseItem.Price)
		if err != nil {
			return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert old price to float64", err)
		}
		if *req.Price != int64(oldPrice) {
			return errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemArrivedNotAllowToChangePrice)
		}
	}

	return nil
}

func (s *Update) updateExpenseAmount(ctx context.Context, tx *sqlx.Tx, expense dbgen.GetStoreExpenseByIDRow, oldExpenseItem dbgen.GetStoreExpenseItemByIDRow, updatedExpenseItem sqlxRepo.UpdateStoreExpenseItemResponse, expenseID, updaterID int64) error {
	oldExpenseAmount, err := utils.PgNumericToFloat64(expense.Amount)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert old expense amount to float64", err)
	}
	oldExpenseItemPrice, err := utils.PgNumericToFloat64(oldExpenseItem.Price)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert old expense item price to float64", err)
	}
	updatedExpenseItemPrice, err := utils.PgNumericToFloat64(updatedExpenseItem.Price)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert updated expense item price to float64", err)
	}

	newExpenseAmount := oldExpenseAmount - (oldExpenseItemPrice * float64(oldExpenseItem.Quantity)) + (updatedExpenseItemPrice * float64(updatedExpenseItem.Quantity))

	return s.repo.Expense.UpdateStoreExpenseAmountTx(ctx, tx, expenseID, sqlxRepo.UpdateStoreExpenseAmountTxParams{
		Amount:  int64(newExpenseAmount),
		Updater: updaterID,
	})
}

func (s *Update) updateProductStockForArrival(ctx context.Context, tx *sqlx.Tx, oldExpenseItem dbgen.GetStoreExpenseItemByIDRow, req adminExpenseItemModel.UpdateParsedRequest) error {
	productID := oldExpenseItem.ProductID
	quantity := int64(oldExpenseItem.Quantity)

	// if pass new product id, use new product id
	if req.ProductID != nil {
		productID = *req.ProductID
	}

	// if pass new quantity, use new quantity
	if req.Quantity != nil {
		quantity = int64(*req.Quantity)
	}

	return s.updateProductStockByAmount(ctx, tx, productID, quantity)
}

func (s *Update) updateProductStockByAmount(ctx context.Context, tx *sqlx.Tx, productID int64, amount int64) error {
	product, err := s.queries.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	newStock := int64(product.CurrentStock) + amount
	return s.repo.Product.UpdateStoreProductStockTx(ctx, tx, productID, sqlxRepo.UpdateStoreProductStockTxParams{
		Stock: int(newStock),
	})
}
