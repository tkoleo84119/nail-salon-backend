package adminStore

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreService struct {
	queries *dbgen.Queries
}

func NewGetStoreService(queries *dbgen.Queries) *GetStoreService {
	return &GetStoreService{
		queries: queries,
	}
}

func (s *GetStoreService) GetStore(ctx context.Context, storeID string) (*adminStoreModel.GetStoreResponse, error) {
	// Parse store ID
	storeIDInt, err := utils.ParseID(storeID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid store ID", err)
	}

	// Get store information
	store, err := s.queries.GetStoreDetailByID(ctx, storeIDInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get store", err)
	}

	// Build response
	response := &adminStoreModel.GetStoreResponse{
		ID:       utils.FormatID(store.ID),
		Name:     store.Name,
		Address:  store.Address.String,
		Phone:    store.Phone.String,
		IsActive: store.IsActive.Bool,
	}

	return response, nil
}
