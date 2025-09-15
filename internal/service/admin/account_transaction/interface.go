package adminAccountTransaction

import (
	"context"

	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.GetAllParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.GetAllResponse, error)
}

type CreateInterface interface {
	Create(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.CreateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.CreateResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, storeID, accountID, transactionID int64, req adminAccountTransactionModel.UpdateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.UpdateResponse, error)
}

type DeleteInterface interface {
	Delete(ctx context.Context, storeID, accountID int64, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.DeleteResponse, error)
}
