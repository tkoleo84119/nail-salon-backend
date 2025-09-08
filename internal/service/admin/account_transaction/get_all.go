package adminAccountTransaction

import (
	"context"

	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(
	repo *sqlxRepo.Repositories,
) *GetAll {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID, accountID int64, req adminAccountTransactionModel.GetAllParsedRequest, creatorStoreIDs []int64) (*adminAccountTransactionModel.GetAllResponse, error) {
	// Check store access
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// Call repository
	total, items, err := s.repo.AccountTransaction.GetAllAccountTransactionsByFilter(ctx, accountID, sqlxRepo.GetAllAccountTransactionsByFilterParams{
		Limit:  &req.Limit,
		Offset: &req.Offset,
		Sort:   &req.Sort,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responseItems := make([]adminAccountTransactionModel.GetAllItem, len(items))
	for i, item := range items {
		responseItems[i] = adminAccountTransactionModel.GetAllItem{
			ID:              utils.FormatID(item.ID),
			TransactionDate: utils.PgTimestamptzToTimeString(item.TransactionDate),
			Type:            item.Type,
			Amount:          int(utils.PgNumericToFloat64(item.Amount)),
			Balance:         int(utils.PgNumericToFloat64(item.Balance)),
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
