package adminReport

import (
	"context"

	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
)

type GetPerformanceMeInterface interface {
	GetPerformanceMe(ctx context.Context, req adminReportModel.GetPerformanceMeParsedRequest, staffID int64, storeIds []int64) (*adminReportModel.GetPerformanceMeResponse, error)
}

type GetStorePerformanceInterface interface {
	GetStorePerformance(ctx context.Context, storeID int64, req adminReportModel.GetStorePerformanceParsedRequest, staffRole string, storeIDs []int64) (*adminReportModel.GetStorePerformanceResponse, error)
}

type GetStoreExpenseInterface interface {
	GetStoreExpense(ctx context.Context, storeID int64, req adminReportModel.GetStoreExpenseParsedRequest, staffRole string, storeIDs []int64) (*adminReportModel.GetStoreExpenseResponse, error)
}
