package adminExpense

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, storeID int64, req adminExpenseModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminExpenseModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// Check if supplierId exists
	supplierExists, err := s.queries.CheckSupplierExists(ctx, req.SupplierID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check supplier existence", err)
	}
	if !supplierExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNotFound)
	}

	var isReimbursed *bool
	if req.PayerID != nil {
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
		// if there is a payer, is_reimbursed is set to false
		falseValue := false
		isReimbursed = &falseValue
	}

	// Create expense data
	expenseID := utils.GenerateID()

	// Convert type
	amountNumeric, err := utils.Int64ToPgNumeric(req.Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount", err)
	}

	_, err = s.queries.CreateExpense(ctx, dbgen.CreateExpenseParams{
		ID:           expenseID,
		StoreID:      storeID,
		Category:     utils.StringPtrToPgText(&req.Category, false),
		SupplierID:   req.SupplierID,
		Amount:       amountNumeric,
		ExpenseDate:  utils.TimeToPgDate(req.ExpenseDate),
		Note:         utils.StringPtrToPgText(req.Note, true),
		PayerID:      utils.Int64PtrToPgInt8(req.PayerID),
		IsReimbursed: utils.BoolPtrToPgBool(isReimbursed),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create expense", err)
	}

	return &adminExpenseModel.CreateResponse{
		ID: utils.FormatID(expenseID),
	}, nil
}
