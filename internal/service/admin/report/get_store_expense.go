package adminReport

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreExpense struct {
	queries *dbgen.Queries
}

func NewGetStoreExpense(queries *dbgen.Queries) GetStoreExpenseInterface {
	return &GetStoreExpense{
		queries: queries,
	}
}

func (s *GetStoreExpense) GetStoreExpense(ctx context.Context, storeID int64, req adminReportModel.GetStoreExpenseParsedRequest, role string, creatorStoreIDs []int64) (*adminReportModel.GetStoreExpenseResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// validate date range not exceed 3 years
	dateDiff := req.EndDate.Sub(req.StartDate)
	if dateDiff > 3*365*24*time.Hour {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ReportDateRangeExceed3Years)
	}

	// get summary
	summary, err := s.queries.GetExpenseReportSummary(ctx, dbgen.GetExpenseReportSummaryParams{
		StoreID:       storeID,
		ExpenseDate:   utils.TimePtrToPgDate(&req.StartDate),
		ExpenseDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense report summary", err)
	}

	// get category stats
	categoryStats, err := s.queries.GetExpenseReportByCategory(ctx, dbgen.GetExpenseReportByCategoryParams{
		StoreID:       storeID,
		ExpenseDate:   utils.TimePtrToPgDate(&req.StartDate),
		ExpenseDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense report by category", err)
	}

	// get supplier stats
	supplierStats, err := s.queries.GetExpenseReportBySupplier(ctx, dbgen.GetExpenseReportBySupplierParams{
		StoreID:       storeID,
		ExpenseDate:   utils.TimePtrToPgDate(&req.StartDate),
		ExpenseDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense report by supplier", err)
	}

	// get payer stats
	payerStats, err := s.queries.GetExpenseReportByPayer(ctx, dbgen.GetExpenseReportByPayerParams{
		StoreID:       storeID,
		ExpenseDate:   utils.TimePtrToPgDate(&req.StartDate),
		ExpenseDate_2: utils.TimePtrToPgDate(&req.EndDate),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense report by payer", err)
	}

	// convert summary
	totalAmount, err := utils.PgNumericToInt64(summary.TotalAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert total amount", err)
	}

	advanceAmount, err := utils.PgNumericToInt64(summary.AdvanceAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert advance amount", err)
	}
	reimbursedAmount, err := utils.PgNumericToInt64(summary.ReimbursedAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert reimbursed amount", err)
	}
	pendingAmount := advanceAmount - reimbursedAmount

	summaryResponse := adminReportModel.ExpenseReportSummary{
		TotalCount:       int(summary.TotalCount),
		TotalAmount:      totalAmount,
		AdvanceAmount:    advanceAmount,
		ReimbursedAmount: reimbursedAmount,
		PendingAmount:    pendingAmount,
	}

	// convert category stats
	categoryResponse := make([]adminReportModel.CategoryExpenseStat, len(categoryStats))
	for i, stat := range categoryStats {
		amount, err := utils.PgNumericToInt64(stat.Amount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert category amount", err)
		}

		categoryResponse[i] = adminReportModel.CategoryExpenseStat{
			Category: stat.Category,
			Count:    int(stat.Count),
			Amount:   amount,
		}
	}

	// convert supplier stats
	supplierResponse := make([]adminReportModel.SupplierExpenseStat, len(supplierStats))
	for i, stat := range supplierStats {
		amount, err := utils.PgNumericToInt64(stat.Amount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert supplier amount", err)
		}

		supplierResponse[i] = adminReportModel.SupplierExpenseStat{
			SupplierID:   utils.FormatID(stat.SupplierID.Int64),
			SupplierName: utils.PgTextToString(stat.SupplierName),
			Count:        int(stat.Count),
			Amount:       amount,
		}
	}

	// convert payer stats
	payerResponse := make([]adminReportModel.PayerExpenseStat, len(payerStats))
	for i, stat := range payerStats {
		advanceAmount, err := utils.PgNumericToInt64(stat.AdvanceAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert payer advance amount", err)
		}
		reimbursedAmount, err := utils.PgNumericToInt64(stat.ReimbursedAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert payer reimbursed amount", err)
		}
		pendingAmount := advanceAmount - reimbursedAmount

		payerResponse[i] = adminReportModel.PayerExpenseStat{
			PayerID:          utils.FormatID(stat.PayerID.Int64),
			PayerName:        utils.PgTextToString(stat.PayerName),
			AdvanceCount:     int(stat.AdvanceCount),
			AdvanceAmount:    advanceAmount,
			ReimbursedAmount: reimbursedAmount,
			PendingAmount:    pendingAmount,
		}
	}

	return &adminReportModel.GetStoreExpenseResponse{
		StartDate: utils.PgDateToDateString(utils.TimePtrToPgDate(&req.StartDate)),
		EndDate:   utils.PgDateToDateString(utils.TimePtrToPgDate(&req.EndDate)),
		Summary:   summaryResponse,
		Category:  categoryResponse,
		Supplier:  supplierResponse,
		Payer:     payerResponse,
	}, nil
}
