package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreAccessService struct {
	repo *sqlxRepo.Repositories
}

func NewCreateStoreAccessService(repo *sqlxRepo.Repositories) *CreateStoreAccessService {
	return &CreateStoreAccessService{
		repo: repo,
	}
}

// CreateStoreAccess creates store access for a staff member
func (s *CreateStoreAccessService) CreateStoreAccess(ctx context.Context, staffID int64, storeID int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStoreAccessResponse, bool, error) {
	// Check if target staff exists
	targetStaff, err := s.repo.Staff.GetStaffUserByID(ctx, staffID)
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
	if targetStaff.Role == adminStaffModel.RoleSuperAdmin {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store exists and is active
	active := true
	_, err = s.repo.Store.GetStoreByID(ctx, storeID, &active)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
		}
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store", err)
	}

	// Check if creator has access to this store (except SUPER_ADMIN)
	if creatorRole != common.RoleSuperAdmin {
		hasAccess, err := utils.CheckStoreAccess(storeID, creatorStoreIDs)
		if err != nil {
			return nil, false, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
		}
		if !hasAccess {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if access already exists
	exists, err := s.repo.StaffUserStoreAccess.CheckStoreAccessExists(ctx, staffID, storeID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store access", err)
	}

	isNewlyCreated := false
	if !exists {
		// Create new store access
		err = s.repo.StaffUserStoreAccess.CreateStaffUserStoreAccess(ctx, sqlxRepo.CreateStaffUserStoreAccessTxParams{
			StoreID:     storeID,
			StaffUserID: staffID,
		})
		if err != nil {
			return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
		}
		isNewlyCreated = true
	}

	// Get complete store access list for the staff member
	storeAccessList, err := s.repo.StaffUserStoreAccess.GetStaffUserStoreAccessByStaffId(ctx, staffID, &active)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store access", err)
	}

	// Convert to response format
	var storeList []common.Store
	for _, access := range storeAccessList {
		storeList = append(storeList, common.Store{
			ID:   utils.FormatID(access.StoreID),
			Name: access.Name,
		})
	}

	response := &adminStaffModel.CreateStoreAccessResponse{
		StoreList: storeList,
	}

	return response, isNewlyCreated, nil
}
