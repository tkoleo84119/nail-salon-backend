package adminStoreAccess

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreAccessModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store_access"
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

func (s *Get) Get(ctx context.Context, staffID int64) (*adminStoreAccessModel.GetResponse, error) {
	// Verify staff user exists
	_, err := s.queries.GetStaffUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff user", err)
	}

	// Get staff store access
	storeAccessList, err := s.queries.GetAllActiveStoreAccessByStaffId(ctx, staffID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get staff store access", err)
	}

	// Convert to response format
	var items []adminStoreAccessModel.GetStore
	for _, access := range storeAccessList {
		items = append(items, adminStoreAccessModel.GetStore{
			ID:   utils.FormatID(access.StoreID),
			Name: access.StoreName,
		})
	}

	return &adminStoreAccessModel.GetResponse{
		StoreList: items,
	}, nil
}
