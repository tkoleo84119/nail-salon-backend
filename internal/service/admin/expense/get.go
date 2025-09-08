package adminExpense

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID, expenseID int64, creatorStoreIDs []int64) (*adminExpenseModel.GetResponse, error) {
	// Check store access permission
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// Get expense basic information
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

	// Get expense items
	expenseItems, err := s.queries.GetStoreExpenseItemsByExpenseID(ctx, expenseID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense items", err)
	}

	amount, err := utils.PgNumericToInt64(expense.Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount to int64", err)
	}

	// Build response
	response := &adminExpenseModel.GetResponse{
		ID:          utils.FormatID(expense.ID),
		Category:    utils.PgTextToString(expense.Category),
		Amount:      amount,
		ExpenseDate: utils.PgDateToDateString(expense.ExpenseDate),
		Note:        utils.PgTextToString(expense.Note),
		CreatedAt:   utils.PgTimestamptzToTimeString(expense.CreatedAt),
		UpdatedAt:   utils.PgTimestamptzToTimeString(expense.UpdatedAt),
	}

	// Add supplier information if exists
	if expense.SupplierID.Valid {
		response.Supplier = &adminExpenseModel.GetExpenseSupplier{
			ID:   utils.FormatID(expense.SupplierID.Int64),
			Name: expense.SupplierName,
		}
	}

	// Add payer information if exists
	if expense.PayerID.Valid {
		response.Payer = &adminExpenseModel.GetExpensePayer{
			ID:   utils.FormatID(expense.PayerID.Int64),
			Name: expense.PayerName,
		}

		// Add reimbursement information
		isReimbursed := utils.PgBoolToBool(expense.IsReimbursed)
		response.IsReimbursed = &isReimbursed

		if expense.ReimbursedAt.Valid {
			reimbursedAt := utils.PgTimestamptzToTimeString(expense.ReimbursedAt)
			response.ReimbursedAt = &reimbursedAt
		}
	}

	if expense.OtherFee.Valid {
		otherFee, err := utils.PgNumericToInt64(expense.OtherFee)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert other fee to int64", err)
		}

		response.OtherFee = &otherFee
	}

	// Add expense items if exist
	if len(expenseItems) > 0 {
		response.Items = make([]adminExpenseModel.GetExpenseItem, len(expenseItems))
		for i, item := range expenseItems {
			price, err := utils.PgNumericToInt64(item.Price)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to int64", err)
			}

			responseItem := adminExpenseModel.GetExpenseItem{
				ID: utils.FormatID(item.ID),
				Product: adminExpenseModel.GetExpenseItemProduct{
					ID:   utils.FormatID(item.ProductID),
					Name: utils.PgTextToString(item.ProductName),
				},
				Quantity:        int64(item.Quantity),
				Price:           price,
				IsArrived:       utils.PgBoolToBool(item.IsArrived),
				StorageLocation: utils.PgTextToString(item.StorageLocation),
				Note:            utils.PgTextToString(item.Note),
			}

			// Add expiration date if exists
			if item.ExpirationDate.Valid {
				expirationDate := utils.PgDateToDateString(item.ExpirationDate)
				responseItem.ExpirationDate = &expirationDate
			}

			// Add arrival date if exists
			if item.ArrivalDate.Valid {
				arrivalDate := utils.PgDateToDateString(item.ArrivalDate)
				responseItem.ArrivalDate = &arrivalDate
			}

			response.Items[i] = responseItem
		}
	}

	return response, nil
}
