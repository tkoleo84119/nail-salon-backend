package adminStaff

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreUsername struct {
	repo *sqlxRepo.Repositories
}

func NewGetStoreUsername(repo *sqlxRepo.Repositories) GetStoreUsernameInterface {
	return &GetStoreUsername{
		repo: repo,
	}
}

func (s *GetStoreUsername) GetStoreUsername(ctx context.Context, storeID int64, req adminStaffModel.GetStoreUsernameParsedRequest, storeIDs []int64) (*adminStaffModel.GetStoreUsernameResponse, error) {
	if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
		return nil, err
	}

	total, items, err := s.repo.Staff.GetAllStaffUsernameByStoreFilter(ctx, storeID, sqlxRepo.GetAllStaffUsernameByStoreFilterParams{
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff username list", err)
	}

	itemsDTO := make([]adminStaffModel.GetStoreUsernameListItem, len(items))
	for i, item := range items {
		itemsDTO[i] = adminStaffModel.GetStoreUsernameListItem{
			ID:       utils.FormatID(item.ID),
			Username: item.Username,
			IsActive: utils.PgBoolToBool(item.IsActive),
		}
	}

	return &adminStaffModel.GetStoreUsernameResponse{
		Total: total,
		Items: itemsDTO,
	}, nil
}
