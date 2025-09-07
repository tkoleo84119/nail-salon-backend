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

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, db *sqlx.DB) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
		db:      db,
	}
}

func (s *Update) Update(ctx context.Context, storeID, expenseID, expenseItemID int64, req adminExpenseItemModel.UpdateParsedRequest, creatorStoreIDs []int64) (*adminExpenseItemModel.UpdateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
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

	oldExpenseItem, err := s.queries.GetStoreExpenseItemByID(ctx, dbgen.GetStoreExpenseItemByIDParams{
		ID:        expenseItemID,
		ExpenseID: expenseID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense item", err)
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
		oldExpenseAmount := utils.PgNumericToFloat64(expense.Amount)
		newExpenseAmount := oldExpenseAmount - (utils.PgNumericToFloat64(oldExpenseItem.Price) * float64(oldExpenseItem.Quantity)) + (utils.PgNumericToFloat64(updatedExpenseItem.Price) * float64(updatedExpenseItem.Quantity))

		// Update expense amount
		err = s.repo.Expense.UpdateStoreExpenseAmountTx(ctx, tx, expenseID, sqlxRepo.UpdateStoreExpenseAmountTxParams{
			Amount: int64(newExpenseAmount),
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense amount", err)
		}
	}

	// update product stock
	if req.Quantity != nil {
		oldProductID := oldExpenseItem.ProductID
		newProductID := req.ProductID

		// if new productId is not nil and new productId is different from old productId
		if newProductID != nil && *newProductID != oldProductID {
			// decrease old product stock
			oldProduct, err := s.queries.GetProductByID(ctx, oldProductID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get old product", err)
			}
			oldProductNewStock := int64(oldProduct.CurrentStock) - int64(oldExpenseItem.Quantity)
			err = s.repo.Product.UpdateStoreProductStockTx(ctx, tx, oldProductID, sqlxRepo.UpdateStoreProductStockTxParams{
				Stock: int(oldProductNewStock),
			})
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update old product stock", err)
			}

			// increase new product stock
			newPID := *newProductID
			newProduct, err := s.queries.GetProductByID(ctx, newPID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get new product", err)
			}
			newProductNewStock := int64(newProduct.CurrentStock) + int64(updatedExpenseItem.Quantity)
			err = s.repo.Product.UpdateStoreProductStockTx(ctx, tx, newPID, sqlxRepo.UpdateStoreProductStockTxParams{
				Stock: int(newProductNewStock),
			})
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update new product stock", err)
			}
		} else {
			// if new productId is nil or new productId is same as old productId
			// adjust old product stock
			oldProduct, err := s.queries.GetProductByID(ctx, oldProductID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get old product", err)
			}
			oldProductNewStock := int64(oldProduct.CurrentStock) - int64(oldExpenseItem.Quantity) + int64(updatedExpenseItem.Quantity)
			err = s.repo.Product.UpdateStoreProductStockTx(ctx, tx, oldProductID, sqlxRepo.UpdateStoreProductStockTxParams{
				Stock: int(oldProductNewStock),
			})
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update product stock", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminExpenseItemModel.UpdateResponse{
		ID: utils.FormatID(updatedExpenseItem.ID),
	}, nil
}
