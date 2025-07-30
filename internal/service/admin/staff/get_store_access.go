package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStaffStoreAccessService struct {
	queries *dbgen.Queries
}

func NewGetStaffStoreAccessService(queries *dbgen.Queries) *GetStaffStoreAccessService {
	return &GetStaffStoreAccessService{
		queries: queries,
	}
}

func (s *GetStaffStoreAccessService) GetStaffStoreAccess(ctx context.Context, staffID string) (*adminStaffModel.GetStaffStoreAccessResponse, error) {
	// Parse staff ID
	staffUserID, err := utils.ParseID(staffID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "Invalid staff ID", err)
	}

	// Verify staff user exists
	_, err = s.queries.GetStaffUserByID(ctx, staffUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get staff user", err)
	}

	// Get staff store access
	storeAccessList, err := s.queries.GetStaffUserStoreAccess(ctx, staffUserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "Failed to get staff store access", err)
	}

	// Convert to response format
	var items []adminStaffModel.StaffStoreAccessItem
	for _, access := range storeAccessList {
		items = append(items, adminStaffModel.StaffStoreAccessItem{
			StoreID: utils.FormatID(access.StoreID),
			Name:    access.StoreName,
		})
	}

	return &adminStaffModel.GetStaffStoreAccessResponse{
		StoreList: items,
	}, nil
}
