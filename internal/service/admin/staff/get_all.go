package adminStaff

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, req adminStaffModel.GetAllParsedRequest) (*adminStaffModel.GetAllResponse, error) {
	total, items, err := s.repo.Staff.GetAllStaffByFilter(ctx, sqlxRepo.GetAllStaffByFilterParams{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff list", err)
	}

	itemsDTO := make([]adminStaffModel.GetAllStaffListItem, len(items))
	for i, item := range items {
		itemsDTO[i] = adminStaffModel.GetAllStaffListItem{
			ID:        utils.FormatID(item.ID),
			Username:  item.Username,
			Email:     item.Email,
			Role:      item.Role,
			IsActive:  utils.PgBoolToBool(item.IsActive),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	return &adminStaffModel.GetAllResponse{
		Total: total,
		Items: itemsDTO,
	}, nil
}
