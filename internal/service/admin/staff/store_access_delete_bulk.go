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

type DeleteStoreAccessBulkService struct {
	repo *sqlxRepo.Repositories
}

func NewDeleteStoreAccessBulkService(repo *sqlxRepo.Repositories) *DeleteStoreAccessBulkService {
	return &DeleteStoreAccessBulkService{
		repo: repo,
	}
}

// DeleteStoreAccessBulk deletes store access for a staff member
func (s *DeleteStoreAccessBulkService) DeleteStoreAccessBulk(ctx context.Context, targetID int64, storeIDs []int64, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.DeleteStoreAccessBulkResponse, error) {
	// Check if target staff exists
	targetStaff, err := s.repo.Staff.GetStaffUserByID(ctx, targetID)
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
	if targetStaff.Role == adminStaffModel.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// For non-SUPER_ADMIN creators, validate they have access to the stores being removed
	if creatorRole != adminStaffModel.RoleSuperAdmin {
		for _, storeID := range storeIDs {
			hasAccess, err := utils.CheckStoreAccess(storeID, creatorStoreIDs)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to check store access", err)
			}
			if !hasAccess {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
			}
		}
	}

	// Delete store access (ignore if access doesn't exist as per spec)
	err = s.repo.StaffUserStoreAccess.BatchDeleteStaffUserStoreAccess(ctx, targetID, storeIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete store access", err)
	}

	// Get remaining store access for the target staff
	remainingStoreAccess, err := s.repo.StaffUserStoreAccess.GetStaffUserStoreAccessByStaffId(ctx, targetID, nil)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get remaining store access", err)
	}

	// Convert to response format
	var storeList []common.Store
	for _, access := range remainingStoreAccess {
		storeList = append(storeList, common.Store{
			ID:   utils.FormatID(access.StoreID),
			Name: access.Name,
		})
	}

	response := &adminStaffModel.DeleteStoreAccessBulkResponse{
		StoreList: storeList,
	}

	return response, nil
}
