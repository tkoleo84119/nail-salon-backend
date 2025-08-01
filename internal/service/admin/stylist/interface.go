package adminStylist

import (
	"context"

	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
)

type UpdateMyStylistServiceInterface interface {
	UpdateMyStylist(ctx context.Context, req adminStylistModel.UpdateMyStylistRequest, staffUserID int64) (*adminStylistModel.UpdateMyStylistResponse, error)
}
