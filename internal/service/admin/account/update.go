package adminAccount

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	accountRepo *sqlx.AccountRepository
}

func NewUpdate(accountRepo *sqlx.AccountRepository) UpdateInterface {
	return &Update{
		accountRepo: accountRepo,
	}
}

func (s *Update) Update(ctx context.Context, accountID int64, req adminAccountModel.UpdateParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountModel.UpdateResponse, error) {
	// Check if account exists and validate store access
	account, err := s.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AccountNotFound, "帳戶不存在或已被刪除", err)
	}

	if err := utils.CheckStoreAccess(account.StoreID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// Update account
	_, err = s.accountRepo.UpdateAccount(ctx, accountID, sqlx.UpdateAccountParams{
		Name:     req.Name,
		Note:     req.Note,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update account", err)
	}

	return &adminAccountModel.UpdateResponse{
		ID: utils.FormatID(accountID),
	}, nil
}
