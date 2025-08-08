package adminStylist

import (
	"context"

	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
)

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, req adminStylistModel.GetAllParsedRequest, role string, storeIds []int64) (*adminStylistModel.GetAllResponse, error)
}

type UpdateMeInterface interface {
	UpdateMe(ctx context.Context, req adminStylistModel.UpdateMeRequest, staffUserID int64) (*adminStylistModel.UpdateMeResponse, error)
}
