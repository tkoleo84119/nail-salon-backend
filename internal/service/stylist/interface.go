package stylist

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
)

type CreateMyStylistServiceInterface interface {
	CreateMyStylist(ctx context.Context, req stylist.CreateMyStylistRequest, staffUserID int64) (*stylist.CreateMyStylistResponse, error)
}

type UpdateMyStylistServiceInterface interface {
	UpdateMyStylist(ctx context.Context, req stylist.UpdateMyStylistRequest, staffUserID int64) (*stylist.UpdateMyStylistResponse, error)
}
