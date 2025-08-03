package adminService

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetServiceService struct {
	repo *sqlxRepo.Repositories
}

func NewGetServiceService(repo *sqlxRepo.Repositories) *GetServiceService {
	return &GetServiceService{
		repo: repo,
	}
}

func (s *GetServiceService) GetService(ctx context.Context, serviceID int64) (*adminServiceModel.GetServiceResponse, error) {
	// Get service information
	service, err := s.repo.Service.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get service", err)
	}

	// Build response
	response := &adminServiceModel.GetServiceResponse{
		ID:              utils.FormatID(service.ID),
		Name:            service.Name,
		DurationMinutes: service.DurationMinutes,
		Price:           int64(utils.PgNumericToFloat64(service.Price)),
		IsAddon:         utils.PgBoolToBool(service.IsAddon),
		IsActive:        utils.PgBoolToBool(service.IsActive),
		IsVisible:       utils.PgBoolToBool(service.IsVisible),
		Note:            service.Note.String,
		CreatedAt:       utils.PgTimestamptzToTimeString(service.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(service.UpdatedAt),
	}

	return response, nil
}
