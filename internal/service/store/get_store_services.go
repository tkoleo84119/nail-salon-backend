package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreServicesService struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
}

func NewGetStoreServicesService(queries dbgen.Querier, repo *sqlxRepo.Repositories) GetStoreServicesServiceInterface {
	return &GetStoreServicesService{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetStoreServicesService) GetStoreServices(ctx context.Context, storeIDStr string, queryParams storeModel.GetStoreServicesQueryParams) (*storeModel.GetStoreServicesResponse, error) {
	// Parse store ID
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid store ID", err)
	}

	// Validate store exists and is active (as per spec requirement)
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
	}

	// Check if store is active
	if !store.IsActive.Bool || !store.IsActive.Valid {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	// Get services from repository with flexible filtering
	services, total, err := s.repo.Service.GetStoreServices(ctx, storeID, queryParams.IsAddon, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store services", err)
	}

	return &storeModel.GetStoreServicesResponse{
		Total: total,
		Items: services,
	}, nil
}
