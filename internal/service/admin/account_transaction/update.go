package adminAccountTransaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
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

func (s *Update) Update(ctx context.Context, storeID, accountID, transactionID int64, req adminAccountTransactionModel.UpdateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.UpdateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	accountTransaction, err := s.queries.GetAccountTransactionByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountTransactionNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get account", err)
	}
	if accountTransaction.AccountID != accountID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AccountTransactionNotBelongToAccount)
	}

	updateResponse, err := s.repo.AccountTransaction.UpdateAccountTransaction(ctx, transactionID, sqlxRepo.UpdateAccountTransactionParams{
		Note: req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update account transaction", err)
	}

	return &adminAccountTransactionModel.UpdateResponse{
		ID: utils.FormatID(updateResponse.ID),
	}, nil
}
