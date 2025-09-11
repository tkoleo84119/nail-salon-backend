package adminExpense

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
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

func (s *Create) Create(ctx context.Context, storeID int64, req adminExpenseModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminExpenseModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// Check if supplierId exists
	if req.SupplierID != nil {
		supplierExists, err := s.queries.CheckSupplierExists(ctx, *req.SupplierID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check supplier existence", err)
		}
		if !supplierExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNotFound)
		}
	}

	var isReimbursed *bool
	if req.PayerID != nil {
		payerHasAccess, err := s.queries.CheckStaffHasStoreAccess(ctx, dbgen.CheckStaffHasStoreAccessParams{
			StaffUserID: *req.PayerID,
			StoreID:     storeID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check payer store access", err)
		}
		if !payerHasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		// if there is a payer, is_reimbursed is set to false
		falseValue := false
		isReimbursed = &falseValue
	}

	expenseID, itemRows, updateProductStockRows, err := s.checkAndPrepareBatchData(ctx, storeID, req.Items)
	if err != nil {
		return nil, err
	}

	// Convert type
	amountNumeric, err := utils.Int64PtrToPgNumeric(&req.Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount", err)
	}
	otherFeeNumeric, err := utils.Int64PtrToPgNumeric(req.OtherFee)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert other fee", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	_, err = qtx.CreateExpense(ctx, dbgen.CreateExpenseParams{
		ID:           expenseID,
		StoreID:      storeID,
		Category:     utils.StringPtrToPgText(&req.Category, false),
		SupplierID:   utils.Int64PtrToPgInt8(req.SupplierID),
		Amount:       amountNumeric,
		OtherFee:     otherFeeNumeric,
		ExpenseDate:  utils.TimePtrToPgDate(&req.ExpenseDate),
		Note:         utils.StringPtrToPgText(req.Note, true),
		PayerID:      utils.Int64PtrToPgInt8(req.PayerID),
		IsReimbursed: utils.BoolPtrToPgBool(isReimbursed),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create expense", err)
	}

	_, err = qtx.BatchCreateExpenseItems(ctx, itemRows)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create expense items", err)
	}

	for _, updateProductStockRow := range updateProductStockRows {
		err = qtx.UpdateProductCurrentStock(ctx, updateProductStockRow)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update product current stock", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminExpenseModel.CreateResponse{
		ID: utils.FormatID(expenseID),
	}, nil
}

func (s *Create) checkAndPrepareBatchData(ctx context.Context, storeID int64, items []adminExpenseModel.CreateItemParsedRequest) (expenseID int64, itemRows []dbgen.BatchCreateExpenseItemsParams, updateProductStockRows []dbgen.UpdateProductCurrentStockParams, err error) {
	productIDs := make([]int64, 0)
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.queries.GetProductsStockInfoByIDs(ctx, productIDs)
	if err != nil {
		return 0, nil, nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get products stock info by ids", err)
	}

	productMap := make(map[int64]dbgen.GetProductsStockInfoByIDsRow)
	for _, product := range products {
		productMap[product.ID] = product

		if product.StoreID != storeID {
			return 0, nil, nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
		}
	}

	expenseID = utils.GenerateID()
	itemRows = make([]dbgen.BatchCreateExpenseItemsParams, 0)
	updateProductStockRows = make([]dbgen.UpdateProductCurrentStockParams, 0)
	now := time.Now()
	nowPg := utils.TimePtrToPgTimestamptz(&now)
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return 0, nil, nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}

		quantity := int32(item.Quantity)

		priceNumeric, err := utils.Int64PtrToPgNumeric(&item.Price)
		if err != nil {
			return 0, nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert total price", err)
		}

		itemRows = append(itemRows, dbgen.BatchCreateExpenseItemsParams{
			ID:              utils.GenerateID(),
			ExpenseID:       expenseID,
			ProductID:       item.ProductID,
			Quantity:        quantity,
			Price:           priceNumeric,
			ExpirationDate:  utils.TimePtrToPgDate(item.ExpirationDate),
			IsArrived:       utils.BoolPtrToPgBool(&item.IsArrived),
			ArrivalDate:     utils.TimePtrToPgDate(item.ArrivalDate),
			StorageLocation: utils.StringPtrToPgText(item.StorageLocation, true),
			Note:            utils.StringPtrToPgText(item.Note, true),
			CreatedAt:       nowPg,
			UpdatedAt:       nowPg,
		})

		updateProductStockRows = append(updateProductStockRows, dbgen.UpdateProductCurrentStockParams{
			ID:           item.ProductID,
			CurrentStock: product.CurrentStock + int32(item.Quantity),
		})
	}

	return expenseID, itemRows, updateProductStockRows, nil
}
