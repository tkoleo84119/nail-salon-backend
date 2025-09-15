package adminExpenseItem

import (
	"context"

	adminExpenseItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense_item"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID, expenseID int64, req adminExpenseItemModel.CreateParsedRequest, creatorID int64, role string, creatorStoreIDs []int64) (*adminExpenseItemModel.CreateResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, storeID, expenseID, expenseItemID int64, req adminExpenseItemModel.UpdateParsedRequest, updaterID int64, role string, creatorStoreIDs []int64) (*adminExpenseItemModel.UpdateResponse, error)
}
