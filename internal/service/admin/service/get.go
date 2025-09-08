package adminService

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, serviceID int64) (*adminServiceModel.GetResponse, error) {
	// Get service information
	service, err := s.queries.GetServiceByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ServiceNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get service", err)
	}

	price, err := utils.PgNumericToInt64(service.Price)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to int64", err)
	}

	// Build response
	response := &adminServiceModel.GetResponse{
		ID:              utils.FormatID(service.ID),
		Name:            service.Name,
		DurationMinutes: service.DurationMinutes,
		Price:           price,
		IsAddon:         utils.PgBoolToBool(service.IsAddon),
		IsActive:        utils.PgBoolToBool(service.IsActive),
		IsVisible:       utils.PgBoolToBool(service.IsVisible),
		Note:            utils.PgTextToString(service.Note),
		CreatedAt:       utils.PgTimestamptzToTimeString(service.CreatedAt),
		UpdatedAt:       utils.PgTimestamptzToTimeString(service.UpdatedAt),
	}

	return response, nil
}
