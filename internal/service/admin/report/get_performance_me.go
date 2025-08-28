package adminReport

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"

	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetPerformanceMe struct {
	queries *dbgen.Queries
}

func NewGetPerformanceMe(queries *dbgen.Queries) *GetPerformanceMe {
	return &GetPerformanceMe{
		queries: queries,
	}
}

func (s *GetPerformanceMe) GetPerformanceMe(ctx context.Context, req adminReportModel.GetPerformanceMeParsedRequest, staffID int64, storeIds []int64) (*adminReportModel.GetPerformanceMeResponse, error) {
	// Validate date range (max 1 year)
	dateDiff := req.EndDate.Sub(req.StartDate)
	if dateDiff > 365*24*time.Hour {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ReportDateRangeExceed1Year)
	}

	stylistID, err := s.queries.GetStylistIDByStaffUserID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StylistNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stylist ID", err)
	}

	// Get performance data by store
	storePerformances, err := s.queries.GetStylistPerformanceGroupByStore(ctx, dbgen.GetStylistPerformanceGroupByStoreParams{
		StylistID:  stylistID,
		WorkDate:   utils.TimeToPgDate(req.StartDate),
		WorkDate_2: utils.TimeToPgDate(req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get performance data", err)
	}

	// If no data found, return empty response
	if len(storePerformances) == 0 {
		return &adminReportModel.GetPerformanceMeResponse{
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
			Stores:            []adminReportModel.GetPerformanceMeStore{},
		}, nil
	}

	// Calculate totals
	var totalBookings, completedBookings, cancelledBookings, noShowBookings, totalServiceTime int
	var totalLinePayRevenue, totalCashRevenue, totalPaidAmount float64

	stores := make([]adminReportModel.GetPerformanceMeStore, len(storePerformances))
	for i, store := range storePerformances {
		// Calculate store totals
		storeTotalBookings := store.CompletedBookings + store.CancelledBookings + store.NoShowBookings
		storeTotalLinePayRevenue := store.LinePayRevenue
		storeTotalCashRevenue := store.CashRevenue
		storeTotalPaidAmount := store.TotalPaidAmount
		storeTotalServiceTime := store.TotalServiceTime

		stores[i] = adminReportModel.GetPerformanceMeStore{
			StoreID:           utils.FormatID(store.StoreID),
			StoreName:         store.StoreName,
			TotalBookings:     int(storeTotalBookings),
			CompletedBookings: int(store.CompletedBookings),
			CancelledBookings: int(store.CancelledBookings),
			NoShowBookings:    int(store.NoShowBookings),
			LinePayRevenue:    utils.PgNumericToFloat64(storeTotalLinePayRevenue),
			CashRevenue:       utils.PgNumericToFloat64(storeTotalCashRevenue),
			TotalPaidAmount:   utils.PgNumericToFloat64(storeTotalPaidAmount),
			TotalServiceTime:  int(storeTotalServiceTime),
		}

		// Accumulate totals
		totalBookings += int(storeTotalBookings)
		completedBookings += int(store.CompletedBookings)
		cancelledBookings += int(store.CancelledBookings)
		noShowBookings += int(store.NoShowBookings)
		totalLinePayRevenue += utils.PgNumericToFloat64(storeTotalLinePayRevenue)
		totalCashRevenue += utils.PgNumericToFloat64(storeTotalCashRevenue)
		totalPaidAmount += utils.PgNumericToFloat64(storeTotalPaidAmount)
		totalServiceTime += int(storeTotalServiceTime)
	}

	return &adminReportModel.GetPerformanceMeResponse{
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
		Stores:            stores,
	}, nil
}
