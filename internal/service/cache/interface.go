package cache

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type AuthCacheInterface interface {
	GetStaffContext(ctx context.Context, userID int64) (*common.StaffContext, error)
	SetStaffContext(ctx context.Context, userID int64, staffContext *common.StaffContext) error
	DeleteStaffContext(ctx context.Context, userID int64) error

	GetCustomerContext(ctx context.Context, customerID int64) (*common.CustomerContext, error)
	SetCustomerContext(ctx context.Context, customerID int64, customerContext *common.CustomerContext) error
	DeleteCustomerContext(ctx context.Context, customerID int64) error
}
