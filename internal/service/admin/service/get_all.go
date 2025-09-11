package adminService

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
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

func (s *GetAll) GetAll(ctx context.Context, req adminServiceModel.GetAllParsedRequest) (*adminServiceModel.GetAllResponse, error) {
	// Get service list from repository
	total, results, err := s.repo.Service.GetAllServicesByFilter(ctx, sqlxRepo.GetAllServicesByFilterParams{
		Name:      req.Name,
		IsAddon:   req.IsAddon,
		IsActive:  req.IsActive,
		IsVisible: req.IsVisible,
		Limit:     &req.Limit,
		Offset:    &req.Offset,
		Sort:      &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get service list", err)
	}

	items := make([]adminServiceModel.GetAllServiceListItemDTO, len(results))
	for i, result := range results {
		price, err := utils.PgNumericToInt64(result.Price)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to int64", err)
		}

		items[i] = adminServiceModel.GetAllServiceListItemDTO{
			ID:              utils.FormatID(result.ID),
			Name:            result.Name,
			Price:           price,
			DurationMinutes: result.DurationMinutes,
			IsAddon:         utils.PgBoolToBool(result.IsAddon),
			IsActive:        utils.PgBoolToBool(result.IsActive),
			IsVisible:       utils.PgBoolToBool(result.IsVisible),
			Note:            utils.PgTextToString(result.Note),
			CreatedAt:       utils.PgTimestamptzToTimeString(result.CreatedAt),
			UpdatedAt:       utils.PgTimestamptzToTimeString(result.UpdatedAt),
		}
	}

	response := &adminServiceModel.GetAllResponse{
		Total: total,
		Items: items,
	}

	return response, nil
}
