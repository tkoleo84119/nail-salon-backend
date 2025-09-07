package adminExpenseItem

import (
	"context"

	adminExpenseItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense_item"
)

type UpdateInterface interface {
	Update(ctx context.Context, storeID, expenseID, expenseItemID int64, req adminExpenseItemModel.UpdateParsedRequest, creatorStoreIDs []int64) (*adminExpenseItemModel.UpdateResponse, error)
}
