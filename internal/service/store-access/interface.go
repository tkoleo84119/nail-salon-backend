package storeAccess

import (
	"context"

	storeAccess "github.com/tkoleo84119/nail-salon-backend/internal/model/store-access"
)

type CreateStoreAccessServiceInterface interface {
	CreateStoreAccess(ctx context.Context, targetID string, req storeAccess.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*storeAccess.CreateStoreAccessResponse, bool, error)
}

type DeleteStoreAccessServiceInterface interface {
	DeleteStoreAccess(ctx context.Context, targetID string, req storeAccess.DeleteStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*storeAccess.DeleteStoreAccessResponse, error)
}
