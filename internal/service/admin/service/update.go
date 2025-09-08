package adminService

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlx.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlx.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, serviceID int64, req adminServiceModel.UpdateRequest, updaterRole string) (*adminServiceModel.UpdateResponse, error) {
	// Check if service exists
	_, err := s.queries.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get service", err)
	}

	// Check name uniqueness if name is being updated
	if req.Name != nil {
		exists, err := s.queries.CheckServiceNameExistsExcluding(ctx, dbgen.CheckServiceNameExistsExcludingParams{
			ID:   serviceID,
			Name: *req.Name,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check service name uniqueness", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceAlreadyExists)
		}
	}

	updatedService, err := s.repo.Service.UpdateService(ctx, serviceID, sqlx.UpdateServiceParams{
		Name:            req.Name,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsAddon:         req.IsAddon,
		IsVisible:       req.IsVisible,
		IsActive:        req.IsActive,
		Note:            req.Note,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update service", err)
	}

	price, err := utils.PgNumericToInt64(updatedService.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to int64", err)
	}

	response := adminServiceModel.UpdateResponse{
		ID:              utils.FormatID(updatedService.ID),
		Name:            updatedService.Name,
		Price:           price,
		DurationMinutes: updatedService.DurationMinutes,
		IsAddon:         utils.PgBoolToBool(updatedService.IsAddon),
		IsVisible:       utils.PgBoolToBool(updatedService.IsVisible),
		IsActive:        utils.PgBoolToBool(updatedService.IsActive),
		Note:            utils.PgTextToString(updatedService.Note),
		CreatedAt:       utils.PgTimestamptzToTimeString(updatedService.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(updatedService.UpdatedAt),
	}

	return &response, nil
}
