package adminService

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetServiceListService struct {
	repo *sqlxRepo.Repositories
}

func NewGetServiceListService(repo *sqlxRepo.Repositories) *GetServiceListService {
	return &GetServiceListService{
		repo: repo,
	}
}

func (s *GetServiceListService) GetServiceList(ctx context.Context, req adminServiceModel.GetServiceListParsedRequest) (*adminServiceModel.GetServiceListResponse, error) {
	// Get service list from repository
	total, results, err := s.repo.Service.GetAllServiceByFilter(ctx, sqlxRepo.GetAllServiceByFilterParams{
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

	items := make([]adminServiceModel.ServiceListItemDTO, len(results))
	for i, result := range results {
		items[i] = adminServiceModel.ServiceListItemDTO{
			ID:              utils.FormatID(result.ID),
			Name:            result.Name,
			Price:           int64(utils.PgNumericToFloat64(result.Price)),
			DurationMinutes: result.DurationMinutes,
			IsAddon:         utils.PgBoolToBool(result.IsAddon),
			IsActive:        utils.PgBoolToBool(result.IsActive),
			IsVisible:       utils.PgBoolToBool(result.IsVisible),
			Note:            utils.PgTextToString(result.Note),
			CreatedAt:       utils.PgTimestamptzToTimeString(result.CreatedAt),
			UpdatedAt:       utils.PgTimestamptzToTimeString(result.UpdatedAt),
		}
	}

	response := &adminServiceModel.GetServiceListResponse{
		Total: total,
		Items: items,
	}

	return response, nil
}
