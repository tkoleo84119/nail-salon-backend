package adminAccountTransaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries *dbgen.Queries
}

func NewDelete(queries *dbgen.Queries) DeleteInterface {
	return &Delete{
		queries: queries,
	}
}

func (s *Delete) Delete(ctx context.Context, storeID, accountID int64, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.DeleteResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// validate account exists
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

	// delete latest account_transactions data
	deletedID, err := s.queries.DeleteLatestAccountTransaction(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountTransactionNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete latest account transaction", err)
	}

	return &adminAccountTransactionModel.DeleteResponse{
		Deleted: utils.FormatID(deletedID),
	}, nil
}
