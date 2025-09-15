package adminAccount

import (
	"context"

	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminAccountModel.CreateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountModel.CreateResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminAccountModel.GetAllParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountModel.GetAllResponse, error)
}

type UpdateInterface interface {
	Update(ctx context.Context, accountID int64, req adminAccountModel.UpdateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountModel.UpdateResponse, error)
}
