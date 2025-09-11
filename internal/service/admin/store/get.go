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

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) GetInterface {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID int64) (*adminStoreModel.GetResponse, error) {
	// Get store information
	store, err := s.queries.GetStoreDetailByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get store", err)
	}

	// Build response
	response := &adminStoreModel.GetResponse{
		ID:        utils.FormatID(store.ID),
		Name:      store.Name,
		Address:   utils.PgTextToString(store.Address),
		Phone:     utils.PgTextToString(store.Phone),
		IsActive:  utils.PgBoolToBool(store.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(store.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(store.UpdatedAt),
	}

	return response, nil
}
