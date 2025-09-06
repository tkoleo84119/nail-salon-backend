package adminExpense

import (
	"context"

	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
)

type CreateInterface interface {
	Create(ctx context.Context, storeID int64, req adminExpenseModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminExpenseModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminExpenseModel.GetAllParsedRequest, creatorStoreIDs []int64) (*adminExpenseModel.GetAllResponse, error)
}

type GetInterface interface {
	Get(ctx context.Context, storeID, expenseID int64, creatorStoreIDs []int64) (*adminExpenseModel.GetResponse, error)
}
