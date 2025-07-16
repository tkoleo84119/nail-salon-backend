package stylist

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
)

type CreateStylistServiceInterface interface {
	CreateStylist(ctx context.Context, req stylist.CreateStylistRequest, staffUserID int64) (*stylist.CreateStylistResponse, error)
}