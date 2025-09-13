package adminCheckout

import (
	"context"

	adminCheckoutModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/checkout"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type CreateBulkInterface interface {
	CreateBulk(ctx context.Context, storeID int64, req adminCheckoutModel.CreateBulkParsedRequest, staffContext *common.StaffContext) (*adminCheckoutModel.CreateBulkResponse, error)
}
