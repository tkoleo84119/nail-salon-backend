package adminAccountTransaction

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(
	repo *sqlxRepo.Repositories,
) GetAllInterface {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.GetAllParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountTransactionModel.GetAllResponse, error) {
	// Check store access
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// Call repository
	total, items, err := s.repo.AccountTransaction.GetAllAccountTransactionsByFilter(ctx, accountID, sqlxRepo.GetAllAccountTransactionsByFilterParams{
		Limit:  &req.Limit,
		Offset: &req.Offset,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responseItems := make([]adminAccountTransactionModel.GetAllItem, len(items))
	for i, item := range items {
		amount, err := utils.PgNumericToInt64(item.Amount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert amount to int64", err)
		}

		balance, err := utils.PgNumericToInt64(item.Balance)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert balance to int64", err)
		}

		responseItems[i] = adminAccountTransactionModel.GetAllItem{
			ID:              utils.FormatID(item.ID),
			TransactionDate: utils.PgDateToDateString(item.TransactionDate),
			Type:            item.Type,
			Amount:          amount,
			Balance:         balance,
			Note:            utils.PgTextToString(item.Note),
			CreatedAt:       utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt:       utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	return &adminAccountTransactionModel.GetAllResponse{
		Total: total,
		Items: responseItems,
	}, nil
}
