package adminAccount

import (
	"context"

	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
)

type CreateInterface interface {
	Create(ctx context.Context, req adminAccountModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminAccountModel.CreateResponse, error)
}
