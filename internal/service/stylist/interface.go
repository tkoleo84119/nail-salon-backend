package stylist

import (
	"context"

	stylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, queryParams stylistModel.GetAllParsedRequest) (*stylistModel.GetAllResponse, error)
}
