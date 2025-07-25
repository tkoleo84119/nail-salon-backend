package adminStylist

import (
	"context"

	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
)

type CreateMyStylistServiceInterface interface {
	CreateMyStylist(ctx context.Context, req adminStylistModel.CreateMyStylistRequest, staffUserID int64) (*adminStylistModel.CreateMyStylistResponse, error)
}

type UpdateMyStylistServiceInterface interface {
	UpdateMyStylist(ctx context.Context, req adminStylistModel.UpdateMyStylistRequest, staffUserID int64) (*adminStylistModel.UpdateMyStylistResponse, error)
}
