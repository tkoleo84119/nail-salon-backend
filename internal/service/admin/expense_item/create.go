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

type Create struct {
	queries *dbgen.Queries
	db      *pgxpool.Pool
}

func NewCreate(queries *dbgen.Queries, db *pgxpool.Pool) CreateInterface {
	return &Create{
		queries: queries,
		db:      db,
	}
}

func (s *Create) Create(ctx context.Context, storeID, expenseID int64, req adminExpenseItemModel.CreateParsedRequest, creatorID int64, role string, creatorStoreIDs []int64) (*adminExpenseItemModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// when pass isArrived is true, but not pass arrivalDate
	if req.IsArrived && req.ArrivalDate == nil {
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

	// if expense is reimbursed, not allow to create expense item
	if expense.IsReimbursed.Valid && expense.IsReimbursed.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseReimbursedNotAllowToCreateItem)
	}

	// if all expense items are arrived, not allow to create expense item
	allArrived, err := s.queries.CheckAllExpenseItemsAreArrived(ctx, expenseID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check all expense items are arrived", err)
	}
	if allArrived {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseItemAllArrivedNotAllowToCreateItem)
	}

	product, err := s.queries.GetProductByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get product", err)
	}
	if product.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
	}

	expenseItemID := utils.GenerateID()
	oldExpenseAmount, err := utils.PgNumericToFloat64(expense.Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert old expense amount to float64", err)
	}
	newExpenseAmount := oldExpenseAmount + (float64(req.Price) * float64(req.Quantity))

	priceNumeric, err := utils.Int64PtrToPgNumeric(&req.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price", err)
	}
	newExpenseAmountNumeric, err := utils.Float64PtrToPgNumeric(&newExpenseAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert new expense amount", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	err = qtx.CreateStoreExpenseItem(ctx, dbgen.CreateStoreExpenseItemParams{
		ID:              expenseItemID,
		ExpenseID:       expenseID,
		ProductID:       req.ProductID,
		Quantity:        int32(req.Quantity),
		Price:           priceNumeric,
		ExpirationDate:  utils.TimePtrToPgDate(req.ExpirationDate),
		IsArrived:       utils.BoolPtrToPgBool(&req.IsArrived),
		ArrivalDate:     utils.TimePtrToPgDate(req.ArrivalDate),
		StorageLocation: utils.StringPtrToPgText(req.StorageLocation, true),
		Note:            utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create expense item", err)
	}

	err = qtx.UpdateStoreExpenseAmount(ctx, dbgen.UpdateStoreExpenseAmountParams{
		Amount:  newExpenseAmountNumeric,
		Updater: utils.Int64PtrToPgInt8(&creatorID),
		ID:      expenseID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense amount", err)
	}

	// if isArrived is true, update product current stock
	if req.IsArrived {
		newProductStock := int64(product.CurrentStock) + req.Quantity
		err = qtx.UpdateProductCurrentStock(ctx, dbgen.UpdateProductCurrentStockParams{
			ID:           req.ProductID,
			CurrentStock: int32(newProductStock),
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update product stock", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminExpenseItemModel.CreateResponse{
		ID: utils.FormatID(expenseItemID),
	}, nil
}
