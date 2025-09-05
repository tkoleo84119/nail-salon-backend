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
		item := adminExpenseModel.GetAllExpenseItem{
			ID: utils.FormatID(expense.ID),
			Supplier: adminExpenseModel.GetAllExpenseSupplierItem{
				ID:   utils.FormatID(expense.SupplierID),
				Name: expense.SupplierName,
			},
			Category:    utils.PgTextToString(expense.Category),
			Amount:      int(utils.PgNumericToFloat64(expense.Amount)),
			ExpenseDate: utils.PgDateToDateString(expense.ExpenseDate),
			Note:        utils.PgTextToString(expense.Note),
			CreatedAt:   utils.PgTimestamptzToTimeString(expense.CreatedAt),
			UpdatedAt:   utils.PgTimestamptzToTimeString(expense.UpdatedAt),
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
