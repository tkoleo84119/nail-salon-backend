package adminAccountTransaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) CreateInterface {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminAccountTransactionModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	account, err := s.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get account", err)
	}
	if account.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountNotBelongToStore)
	}

	balance, err := s.queries.GetAccountTransactionCurrentBalance(ctx, accountID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get account transaction current balance", err)
	}

	if req.Type == common.AccountTransactionTypeIncome {
		balance += int32(req.Amount)
	} else {
		balance -= int32(req.Amount)
	}

	// check if balance is less than 0
	if balance < 0 {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountBalanceNotEnough)
	}

	int64Balance := int64(balance)
	balanceNumeric, err := utils.Int64PtrToPgNumeric(&int64Balance)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert balance", err)
	}

	int64Amount := int64(req.Amount)
	amountNumeric, err := utils.Int64PtrToPgNumeric(&int64Amount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount", err)
	}

	// Create account
	accountTransactionID := utils.GenerateID()
	_, err = s.queries.CreateAccountTransaction(ctx, dbgen.CreateAccountTransactionParams{
		ID:              accountTransactionID,
		AccountID:       accountID,
		TransactionDate: utils.TimePtrToPgDate(&req.TransactionDate),
		Type:            req.Type,
		Amount:          amountNumeric,
		Balance:         balanceNumeric,
		Note:            utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create account transaction", err)
	}

	return &adminAccountTransactionModel.CreateResponse{
		ID: utils.FormatID(accountID),
	}, nil
}
