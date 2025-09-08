package adminAccountTransaction

import (
	"context"

	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.GetAllParsedRequest, creatorStoreIDs []int64) (*adminAccountTransactionModel.GetAllResponse, error)
}

type CreateInterface interface {
	Create(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminAccountTransactionModel.CreateResponse, error)
}
