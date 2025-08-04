package adminStylist

import (
	"context"

	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
)

type GetStylistListServiceInterface interface {
	GetStylistList(ctx context.Context, storeID int64, req adminStylistModel.GetStylistListParsedRequest, role string, storeIds []int64) (*adminStylistModel.GetStylistListResponse, error)
}

type UpdateMyStylistServiceInterface interface {
	UpdateMyStylist(ctx context.Context, req adminStylistModel.UpdateMyStylistRequest, staffUserID int64) (*adminStylistModel.UpdateMyStylistResponse, error)
}
