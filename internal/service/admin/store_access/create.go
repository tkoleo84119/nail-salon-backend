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

type Create struct {
	queries   *dbgen.Queries
	authCache cache.AuthCacheInterface
}

func NewCreate(queries *dbgen.Queries, authCache cache.AuthCacheInterface) *Create {
	return &Create{
		queries:   queries,
		authCache: authCache,
	}
}

// CreateStoreAccess creates store access for a staff member
func (s *Create) Create(ctx context.Context, staffID int64, storeID int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStoreAccessModel.CreateResponse, bool, error) {
	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get target staff", err)
	}

	// Cannot modify self
	if staffID == creatorID {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotUpdateSelf)
	}

	// Cannot modify SUPER_ADMIN
	if targetStaff.Role == common.RoleSuperAdmin {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if creator has access to this store (except SUPER_ADMIN)
	if creatorRole != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
			return nil, false, err
		}
	}

	// Check if access already exists
	exists, err := s.queries.CheckStoreAccessExists(ctx, dbgen.CheckStoreAccessExistsParams{
		StaffUserID: staffID,
		StoreID:     storeID,
	})
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store access", err)
	}

	isNewlyCreated := false
	if !exists {
		// Create new store access
		err = s.queries.CreateStaffUserStoreAccess(ctx, dbgen.CreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: staffID,
		})
		if err != nil {
			return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
		}
		isNewlyCreated = true
	}

	// Get complete store access list for the staff member
	storeAccessList, err := s.queries.GetAllActiveStoreAccessByStaffId(ctx, staffID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store access", err)
	}

	// Convert to response format
	var storeList []adminStoreAccessModel.CreateStore
	for _, access := range storeAccessList {
		storeList = append(storeList, adminStoreAccessModel.CreateStore{
			ID:   utils.FormatID(access.StoreID),
			Name: access.StoreName,
		})
	}

	response := &adminStoreAccessModel.CreateResponse{
		StoreList: storeList,
	}

	// if newly created store access, delete staff context from cache
	if isNewlyCreated {
		if cacheErr := s.authCache.DeleteStaffContext(ctx, staffID); cacheErr != nil {
			log.Println("failed to delete staff context from cache", cacheErr)
		}
	}

	return response, isNewlyCreated, nil
}
