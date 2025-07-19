package store

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store"
)

type CreateStoreServiceInterface interface {
	CreateStore(ctx context.Context, req store.CreateStoreRequest, staffContext common.StaffContext) (*store.CreateStoreResponse, error)
}