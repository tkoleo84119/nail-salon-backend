package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreService struct {
	queries dbgen.Querier
}

func NewGetStoreService(queries dbgen.Querier) GetStoreServiceInterface {
	return &GetStoreService{
		queries: queries,
	}
}

func (s *GetStoreService) GetStore(ctx context.Context, storeIDStr string) (*storeModel.GetStoreResponse, error) {
	// Parse store ID
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	// Get store details using existing SQLC query
	store, err := s.queries.GetStoreDetailByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
	}

	// Check if store is active (as per spec requirement)
	if !store.IsActive.Bool || !store.IsActive.Valid {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	// Build response
	return &storeModel.GetStoreResponse{
		ID:      utils.FormatID(store.ID),
		Name:    store.Name,
		Address: utils.PgTextToString(store.Address),
		Phone:   utils.PgTextToString(store.Phone),
	}, nil
}