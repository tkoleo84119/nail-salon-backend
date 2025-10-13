package adminExpense

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID, expenseID int64, req adminExpenseModel.UpdateParsedRequest, updaterID int64, role string, creatorStoreIDs []int64) (*adminExpenseModel.UpdateResponse, error) {
	// Check store access permission
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// when pass isReimbursed is true, but not pass reimbursedAt
	if req.IsReimbursed != nil && *req.IsReimbursed && req.ReimbursedAt == nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotUpdateReimbursedInfoWithoutReimbursedAt)
	}

	// Verify expense exists and belongs to the store
	expense, err := s.queries.GetStoreExpenseByID(ctx, dbgen.GetStoreExpenseByIDParams{
		ID:      expenseID,
		StoreID: storeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get expense", err)
	}

	// if expense is reimbursed, only allow to update note
	if expense.IsReimbursed.Valid && expense.IsReimbursed.Bool {
		if req.IsReimbursed != nil || req.ReimbursedAt != nil || req.PayerID != nil || req.SupplierID != nil || req.ExpenseDate != nil || req.OtherFee != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotUpdateReimbursedInfo)
		}
	}

	// Validate supplier if provided
	if req.SupplierID != nil {
		supplierExists, err := s.queries.CheckSupplierExists(ctx, *req.SupplierID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check supplier existence", err)
		}
		if !supplierExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNotFound)
		}
	}

	// Validate payer and check store access if provided
	if req.PayerID != nil {
		// Check if staff exists and has access to the store
		payerHasAccess, err := s.queries.CheckStaffHasStoreAccess(ctx, dbgen.CheckStaffHasStoreAccessParams{
			StaffUserID: *req.PayerID,
			StoreID:     storeID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check payer store access", err)
		}
		if !payerHasAccess {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
	}

	// Check if isReimbursed or reimbursedAt is provided without payerId
	if (req.IsReimbursed != nil || req.ReimbursedAt != nil) && req.PayerID == nil && !expense.PayerID.Valid {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotUpdateReimbursedInfoWithoutPayerID)
	}

	// if amount or category is being modified and there are expense items exist, not allow to update
	if req.Amount != nil || req.Category != nil {
		expenseItemsExists, err := s.queries.CheckExpenseItemsExistsByExpenseID(ctx, expenseID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check expense items existence", err)
		}
		if expenseItemsExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseNotUpdateAmountOrCategoryWithExpenseItems)
		}
	}

	// if isReimbursed is true, and expense items exist, check all expense items are arrived
	if req.IsReimbursed != nil && *req.IsReimbursed {
		expenseItemsExists, err := s.queries.CheckExpenseItemsExistsByExpenseID(ctx, expenseID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check expense items existence", err)
		}

		if expenseItemsExists {
			allArrived, err := s.queries.CheckAllExpenseItemsAreArrived(ctx, expenseID)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check all expense items are arrived", err)
			}
			if !allArrived {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ExpenseReimbursementNotAllowItemNotArrived)
			}
		}
	}

	// Update expense
	updateParams := sqlxRepo.UpdateStoreExpenseParams{
		SupplierID:   req.SupplierID,
		Category:     req.Category,
		Amount:       req.Amount,
		OtherFee:     req.OtherFee,
		ExpenseDate:  req.ExpenseDate,
		Note:         req.Note,
		PayerID:      req.PayerID,
		IsReimbursed: req.IsReimbursed,
		ReimbursedAt: req.ReimbursedAt,
		Updater:      &updaterID,
	}

	if req.PayerIDIsNone != nil && *req.PayerIDIsNone {
		updateParams.PayerIDIsNone = req.PayerIDIsNone

		// Set payerID, isReimbursed, and reimbursedAt to nil, avoid updating these fields
		updateParams.PayerID = nil
		updateParams.IsReimbursed = nil
		updateParams.ReimbursedAt = nil
	}

	result, err := s.repo.Expense.UpdateStoreExpense(ctx, storeID, expenseID, updateParams)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update expense", err)
	}

	return &adminExpenseModel.UpdateResponse{
		ID: utils.FormatID(result.ID),
	}, nil
}
