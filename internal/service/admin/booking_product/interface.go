package adminBookingProduct

import (
	"context"

	adminBookingProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking_product"
)

type BulkCreateInterface interface {
	BulkCreate(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.BulkCreateParsedRequest, staffStoreIDs []int64) (*adminBookingProductModel.BulkCreateResponse, error)
}

type BulkDeleteInterface interface {
	BulkDelete(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.BulkDeleteParsedRequest, staffStoreIDs []int64) (*adminBookingProductModel.BulkDeleteResponse, error)
}

type GetAllInterface interface {
	GetAll(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.GetAllParsedRequest, staffStoreIDs []int64) (*adminBookingProductModel.GetAllResponse, error)
}
