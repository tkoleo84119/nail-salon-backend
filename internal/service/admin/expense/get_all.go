package adminExpense

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminExpenseModel.GetAllParsedRequest, creatorStoreIDs []int64) (*adminExpenseModel.GetAllResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	total, expenses, err := s.repo.Expense.GetAllStoreExpensesByFilter(ctx, storeID, sqlxRepo.GetAllStoreExpensesByFilterParams{
		Category:     req.Category,
		SupplierID:   req.SupplierID,
		PayerID:      req.PayerID,
		IsReimbursed: req.IsReimbursed,
		Limit:        &req.Limit,
		Offset:       &req.Offset,
		Sort:         &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expenses", err)
	}

	items := make([]adminExpenseModel.GetAllExpenseItem, len(expenses))
	for i, expense := range expenses {
		amount, err := utils.PgNumericToInt64(expense.Amount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount to int64", err)
		}

		item := adminExpenseModel.GetAllExpenseItem{
			ID:          utils.FormatID(expense.ID),
			Category:    utils.PgTextToString(expense.Category),
			Amount:      amount,
			ExpenseDate: utils.PgDateToDateString(expense.ExpenseDate),
			Note:        utils.PgTextToString(expense.Note),
			CreatedAt:   utils.PgTimestamptzToTimeString(expense.CreatedAt),
			UpdatedAt:   utils.PgTimestamptzToTimeString(expense.UpdatedAt),
		}

		if expense.SupplierID.Valid {
			item.Supplier = &adminExpenseModel.GetAllExpenseSupplierItem{
				ID:   utils.FormatID(expense.SupplierID.Int64),
				Name: expense.SupplierName,
			}
		}

		if expense.OtherFee.Valid {
			otherFee, err := utils.PgNumericToInt64(expense.OtherFee)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert other fee to int64", err)
			}

			item.OtherFee = &otherFee
		}

		// Only include payer, isReimbursed, and reimbursedAt if payer exists
		if expense.PayerID.Valid {
			item.Payer = &adminExpenseModel.GetAllExpensePayerItem{
				ID:   utils.FormatID(expense.PayerID.Int64),
				Name: utils.PgTextToString(expense.PayerName),
			}
			isReimbursed := utils.PgBoolToBool(expense.IsReimbursed)
			item.IsReimbursed = &isReimbursed

			if expense.ReimbursedAt.Valid {
				reimbursedAt := utils.PgTimestamptzToTimeString(expense.ReimbursedAt)
				item.ReimbursedAt = &reimbursedAt
			}
		}

		items[i] = item
	}

	return &adminExpenseModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
