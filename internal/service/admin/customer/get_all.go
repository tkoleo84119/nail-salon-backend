package adminCustomer

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
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

func (s *GetAll) GetAll(ctx context.Context, req adminCustomerModel.GetAllParsedRequest) (*adminCustomerModel.GetAllResponse, error) {
	// Get customers from repository
	total, results, err := s.repo.Customer.GetAllCustomersByFilter(ctx, sqlxRepo.GetAllCustomersByFilterParams{
		Name:          req.Name,
		Phone:         req.Phone,
		Level:         req.Level,
		IsBlacklisted: req.IsBlacklisted,
		MinPastDays:   req.MinPastDays,
		Limit:         &req.Limit,
		Offset:        &req.Offset,
		Sort:          &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get customers", err)
	}

	items := make([]adminCustomerModel.GetAllCustomerItem, len(results))
	for i, result := range results {
		items[i] = adminCustomerModel.GetAllCustomerItem{
			ID:            utils.FormatID(result.ID),
			Name:          result.Name,
			Phone:         result.Phone,
			Birthday:      utils.PgDateToDateString(result.Birthday),
			City:          utils.PgTextToString(result.City),
			Level:         utils.PgTextToString(result.Level),
			IsBlacklisted: utils.PgBoolToBool(result.IsBlacklisted),
			LastVisitAt:   utils.PgTimestamptzToTimeString(result.LastVisitAt),
			UpdatedAt:     utils.PgTimestamptzToTimeString(result.UpdatedAt),
		}
	}

	return &adminCustomerModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
