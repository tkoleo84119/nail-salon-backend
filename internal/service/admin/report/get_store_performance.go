package adminReport

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStorePerformance struct {
	queries *dbgen.Queries
}

func NewGetStorePerformance(queries *dbgen.Queries) GetStorePerformanceInterface {
	return &GetStorePerformance{
		queries: queries,
	}
}

func (s *GetStorePerformance) GetStorePerformance(ctx context.Context, storeID int64, req adminReportModel.GetStorePerformanceParsedRequest, staffRole string, storeIDs []int64) (*adminReportModel.GetStorePerformanceResponse, error) {
	if staffRole != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
			return nil, err
		}
	}

	// Validate date range (max 1 year)
	dateDiff := req.EndDate.Sub(req.StartDate)
	if dateDiff > 365*24*time.Hour {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ReportDateRangeExceed1Year)
	}

	// Get performance data by stylist for the specific store
	stylistPerformances, err := s.queries.GetStorePerformanceGroupByStylist(ctx, dbgen.GetStorePerformanceGroupByStylistParams{
		StoreID:    storeID,
		WorkDate:   utils.TimePtrToPgDate(&req.StartDate),
		WorkDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store performance data", err)
	}

	// if no data found, return empty response
	if len(stylistPerformances) == 0 {
		return &adminReportModel.GetStorePerformanceResponse{
			StartDate:         req.StartDate.Format("2006-01-02"),
			EndDate:           req.EndDate.Format("2006-01-02"),
			TotalBookings:     0,
			CompletedBookings: 0,
			CancelledBookings: 0,
			NoShowBookings:    0,
			LinePayRevenue:    0,
			CashRevenue:       0,
			TotalPaidAmount:   0,
			TotalServiceTime:  0,
			Stylists:          make([]adminReportModel.GetStorePerformanceStylist, 0),
		}, nil
	}

	var totalBookings, completedBookings, cancelledBookings, noShowBookings, totalServiceTime int
	var totalLinePayRevenue, totalCashRevenue, totalPaidAmount float64

	stylists := make([]adminReportModel.GetStorePerformanceStylist, len(stylistPerformances))
	// Process stylist performances and calculate totals
	for i, stylist := range stylistPerformances {
		// stylist totals
		stylistTotalBookings := stylist.CompletedBookings + stylist.CancelledBookings + stylist.NoShowBookings
		stylistTotalLinePayRevenue := stylist.LinePayRevenue
		stylistTotalCashRevenue := stylist.CashRevenue
		stylistTotalPaidAmount := stylist.TotalPaidAmount
		stylistTotalServiceTime := stylist.TotalServiceTime

		stylistTotalLinePayRevenueFloat, err := utils.PgNumericToFloat64(stylistTotalLinePayRevenue)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert stylist total line pay revenue to float64", err)
		}
		stylistTotalCashRevenueFloat, err := utils.PgNumericToFloat64(stylistTotalCashRevenue)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert stylist total cash revenue to float64", err)
		}
		stylistTotalPaidAmountFloat, err := utils.PgNumericToFloat64(stylistTotalPaidAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert stylist total paid amount to float64", err)
		}

		stylists[i] = adminReportModel.GetStorePerformanceStylist{
			StylistID:         utils.FormatID(stylist.StylistID),
			StylistName:       utils.PgTextToString(stylist.StylistName),
			TotalBookings:     int(stylistTotalBookings),
			CompletedBookings: int(stylist.CompletedBookings),
			CancelledBookings: int(stylist.CancelledBookings),
			NoShowBookings:    int(stylist.NoShowBookings),
			LinePayRevenue:    stylistTotalLinePayRevenueFloat,
			CashRevenue:       stylistTotalCashRevenueFloat,
			TotalPaidAmount:   stylistTotalPaidAmountFloat,
			TotalServiceTime:  int(stylistTotalServiceTime),
		}

		// Add to totals
		totalBookings += int(stylistTotalBookings)
		completedBookings += int(stylist.CompletedBookings)
		cancelledBookings += int(stylist.CancelledBookings)
		noShowBookings += int(stylist.NoShowBookings)
		totalLinePayRevenue += stylistTotalLinePayRevenueFloat
		totalCashRevenue += stylistTotalCashRevenueFloat
		totalPaidAmount += stylistTotalPaidAmountFloat
		totalServiceTime += int(stylistTotalServiceTime)
	}

	return &adminReportModel.GetStorePerformanceResponse{
		StartDate:         req.StartDate.Format("2006-01-02"),
		EndDate:           req.EndDate.Format("2006-01-02"),
		TotalBookings:     totalBookings,
		CompletedBookings: completedBookings,
		CancelledBookings: cancelledBookings,
		NoShowBookings:    noShowBookings,
		LinePayRevenue:    totalLinePayRevenue,
		CashRevenue:       totalCashRevenue,
		TotalPaidAmount:   totalPaidAmount,
		TotalServiceTime:  totalServiceTime,
		Stylists:          stylists,
	}, nil
}
