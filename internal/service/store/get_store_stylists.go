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

type GetStoreStylistsService struct {
	queries     dbgen.Querier
	stylistRepo *sqlxRepo.StylistRepository
}

func NewGetStoreStylistsService(queries dbgen.Querier, stylistRepo *sqlxRepo.StylistRepository) GetStoreStylistsServiceInterface {
	return &GetStoreStylistsService{
		queries:     queries,
		stylistRepo: stylistRepo,
	}
}

func (s *GetStoreStylistsService) GetStoreStylists(ctx context.Context, storeIDStr string, queryParams storeModel.GetStoreStylistsQueryParams) (*storeModel.GetStoreStylistsResponse, error) {
	// Parse store ID
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
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

	// Get stylists from repository
	stylists, total, err := s.stylistRepo.GetStoreStylists(ctx, storeID, queryParams.Limit, queryParams.Offset)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store stylists", err)
	}

	return &storeModel.GetStoreStylistsResponse{
		Total: total,
		Items: stylists,
	}, nil
}
