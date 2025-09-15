package adminStoreAccess

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreAccessModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store_access"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteBulk struct {
	queries   *dbgen.Queries
	authCache cache.AuthCacheInterface
}

func NewDeleteBulk(queries *dbgen.Queries, authCache cache.AuthCacheInterface) DeleteBulkInterface {
	return &DeleteBulk{
		queries:   queries,
		authCache: authCache,
	}
}

// DeleteStoreAccessBulk deletes store access for a staff member
func (s *DeleteBulk) DeleteBulk(ctx context.Context, targetID int64, storeIDs []int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStoreAccessModel.DeleteBulkResponse, error) {
	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get target staff", err)
	}

	// Cannot modify self
	if targetID == creatorID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotUpdateSelf)
	}
	// Cannot modify SUPER_ADMIN store access
	if targetStaff.Role == common.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// For non-SUPER_ADMIN creators, validate they have access to the stores being removed
	for _, storeID := range storeIDs {
		if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, creatorRole); err != nil {
			return nil, err
		}
	}

	// Delete store access
	err = s.queries.DeleteStaffUserStoreAccess(ctx, dbgen.DeleteStaffUserStoreAccessParams{
		StaffUserID: targetID,
		Column2:     storeIDs,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete store access", err)
	}

	// Convert to response format
	var deleted []string
	for _, storeID := range storeIDs {
		deleted = append(deleted, utils.FormatID(storeID))
	}

	response := &adminStoreAccessModel.DeleteBulkResponse{
		Deleted: deleted,
	}

	// delete staff context from cache
	if cacheErr := s.authCache.DeleteStaffContext(ctx, targetID); cacheErr != nil {
		log.Println("failed to delete staff context from cache", cacheErr)
	}

	return response, nil
}
