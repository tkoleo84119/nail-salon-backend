package adminSupplier

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		repo: repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, req adminSupplierModel.GetAllParsedRequest) (*adminSupplierModel.GetAllResponse, error) {
	total, items, err := s.repo.Supplier.GetAllSuppliersByFilter(ctx, sqlxRepo.GetAllSuppliersByFilterParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get supplier list", err)
	}

	itemsDTO := make([]adminSupplierModel.GetAllSupplierItem, len(items))
	for i, item := range items {
		itemsDTO[i] = adminSupplierModel.GetAllSupplierItem{
			ID:        utils.FormatID(item.ID),
			Name:      item.Name,
			IsActive:  utils.PgBoolToBool(item.IsActive),
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(item.UpdatedAt),
		}
	}

	response := &adminSupplierModel.GetAllResponse{
		Total: total,
		Items: itemsDTO,
	}

	return response, nil
}
