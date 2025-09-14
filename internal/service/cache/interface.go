package cache

import (
	"context"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
)

type AuthCacheInterface interface {
	GetStaffContext(ctx context.Context, userID int64) (*common.StaffContext, error)
	SetStaffContext(ctx context.Context, userID int64, staffContext *common.StaffContext) error
	DeleteStaffContext(ctx context.Context, userID int64) error
	DeleteAllStaffContext(ctx context.Context) error

	GetCustomerContext(ctx context.Context, customerID int64) (*common.CustomerContext, error)
	SetCustomerContext(ctx context.Context, customerID int64, customerContext *common.CustomerContext) error
	DeleteCustomerContext(ctx context.Context, customerID int64) error
}

type ActivityLogCacheInterface interface {
	LogActivity(ctx context.Context, activityType common.ActivityLogType, message string) error
	GetRecentActivities(ctx context.Context, limit int) (*common.ActivityLogResponse, error)

	LogCustomerRegister(ctx context.Context, customerName string, lineName string) error
	LogCustomerBrowse(ctx context.Context, customerName string, lineName string) error
	LogCustomerBooking(ctx context.Context, customerName string, lineName string, storeName string) error
	LogCustomerBookingUpdate(ctx context.Context, customerName string, lineName string, storeName string) error
	LogCustomerBookingCancel(ctx context.Context, customerName string, lineName string, storeName string) error
	LogAdminBookingCreate(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error
	LogAdminBookingUpdate(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error
	LogAdminBookingCancel(ctx context.Context, staffName string, customerName string, lineName string, storeName string) error
	LogAdminBookingCompleted(ctx context.Context, staffName string, customerName string, lineName string, checkoutCount int, storeName string) error
}
