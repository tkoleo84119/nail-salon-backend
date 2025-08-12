package service

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	serviceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/service"
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

func (s *GetAll) GetAll(ctx context.Context, queryParams serviceModel.GetAllParsedRequest) (*serviceModel.GetAllResponse, error) {
	trueCondition := true
	visibleCondition := true
	// Get services from repository with flexible filtering
	total, services, err := s.repo.Service.GetAllServiceByFilter(ctx, sqlxRepo.GetAllServiceByFilterParams{
		IsActive:  &trueCondition,
		IsVisible: &visibleCondition,
		IsAddon:   queryParams.IsAddon,
		Limit:     &queryParams.Limit,
		Offset:    &queryParams.Offset,
		Sort:      &queryParams.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store services", err)
	}

	items := make([]serviceModel.GetAllItem, len(services))
	for i, service := range services {
		items[i] = serviceModel.GetAllItem{
			ID:              utils.FormatID(service.ID),
			Name:            service.Name,
			Price:           int(utils.PgNumericToFloat64(service.Price)),
			DurationMinutes: int(service.DurationMinutes),
			IsAddon:         utils.PgBoolToBool(service.IsAddon),
			Note:            utils.PgTextToString(service.Note),
		}
	}

	return &serviceModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
