package adminAccount

import (
	"context"

	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	accountRepo *sqlx.AccountRepository
}

func NewGetAll(accountRepo *sqlx.AccountRepository) GetAllInterface {
	return &GetAll{
		accountRepo: accountRepo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminAccountModel.GetAllParsedRequest, role string, creatorStoreIDs []int64) (*adminAccountModel.GetAllResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	// Call repository
	total, items, err := s.accountRepo.GetAllAccountsByFilter(ctx, sqlx.GetAllAccountsByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responseItems := make([]adminAccountModel.GetAllItem, len(items))
	for i, item := range items {
		responseItems[i] = adminAccountModel.GetAllItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			Note:      utils.PgTextToString(item.Note),
			IsActive:  utils.PgBoolToBool(item.IsActive),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	return &adminAccountModel.GetAllResponse{
		Total: total,
		Items: responseItems,
	}, nil
}
