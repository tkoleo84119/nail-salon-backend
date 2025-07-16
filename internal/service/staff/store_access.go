package staff

import (
	"context"
	"database/sql"
	"fmt"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreAccessService struct {
	queries dbgen.Querier
}

func NewCreateStoreAccessService(queries dbgen.Querier) *CreateStoreAccessService {
	return &CreateStoreAccessService{
		queries: queries,
	}
}

// CreateStoreAccess creates store access for a staff member
func (s *CreateStoreAccessService) CreateStoreAccess(ctx context.Context, targetID string, req staff.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*staff.CreateStoreAccessResponse, bool, error) {
	// Parse target staff ID
	targetStaffID, err := utils.ParseID(targetID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid target staff ID", err)
	}

	// Parse store ID
	storeID, err := utils.ParseID(req.StoreID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
	}

	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.UserStaffNotFound)
		}
		return nil, false, fmt.Errorf("failed to get target staff: %w", err)
	}

	// Cannot modify self
	if targetStaffID == creatorID {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.UserNotUpdateSelf)
	}

	// Cannot modify SUPER_ADMIN
	if targetStaff.Role == staff.RoleSuperAdmin {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store exists and is active
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.UserStoreNotFound)
		}
		return nil, false, fmt.Errorf("failed to get store: %w", err)
	}

	if !store.IsActive.Bool {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.UserStoreNotActive)
	}

	// Check if creator has access to this store (except SUPER_ADMIN)
	if creatorRole != staff.RoleSuperAdmin {
		hasAccess := false
		for _, creatorStoreID := range creatorStoreIDs {
			if creatorStoreID == storeID {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	// Check if access already exists
	existsResult, err := s.queries.CheckStoreAccessExists(ctx, dbgen.CheckStoreAccessExistsParams{
		StaffUserID: targetStaffID,
		StoreID:     storeID,
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to check store access existence: %w", err)
	}

	isNewlyCreated := false
	if !existsResult {
		// Create new store access
		err = s.queries.CreateStaffUserStoreAccess(ctx, dbgen.CreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: targetStaffID,
		})
		if err != nil {
			return nil, false, fmt.Errorf("failed to create store access: %w", err)
		}
		isNewlyCreated = true
	}

	// Get complete store access list for the staff member
	storeAccessList, err := s.queries.GetStaffUserStoreAccess(ctx, targetStaffID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get staff store access: %w", err)
	}

	// Convert to response format
	var storeList []common.Store
	for _, access := range storeAccessList {
		storeList = append(storeList, common.Store{
			ID:   utils.FormatID(access.StoreID),
			Name: access.StoreName,
		})
	}

	response := &staff.CreateStoreAccessResponse{
		StaffUserID: targetID,
		StoreList:   storeList,
	}

	return response, isNewlyCreated, nil
}

type DeleteStoreAccessService struct {
	queries dbgen.Querier
}

func NewDeleteStoreAccessService(queries dbgen.Querier) *DeleteStoreAccessService {
	return &DeleteStoreAccessService{
		queries: queries,
	}
}

// DeleteStoreAccess deletes store access for a staff member
func (s *DeleteStoreAccessService) DeleteStoreAccess(ctx context.Context, targetID string, req staff.DeleteStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*staff.DeleteStoreAccessResponse, error) {
	// Parse target staff ID
	targetStaffID, err := utils.ParseID(targetID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid target staff ID", err)
	}

	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetStaffID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserStaffNotFound)
		}
		return nil, fmt.Errorf("failed to get target staff: %w", err)
	}

	// Cannot modify self
	if targetStaffID == creatorID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserNotUpdateSelf)
	}

	// Cannot modify SUPER_ADMIN store access
	if targetStaff.Role == staff.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Parse store IDs
	var storeIDs []int64
	for _, storeIDStr := range req.StoreIDs {
		storeID, err := utils.ParseID(storeIDStr)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid store ID", err)
		}
		storeIDs = append(storeIDs, storeID)
	}

	// For non-SUPER_ADMIN creators, validate they have access to the stores being removed
	if creatorRole != staff.RoleSuperAdmin {
		for _, storeID := range storeIDs {
			hasAccess := false
			for _, creatorStoreID := range creatorStoreIDs {
				if storeID == creatorStoreID {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
			}
		}
	}

	// Delete store access (ignore if access doesn't exist as per spec)
	err = s.queries.DeleteStaffUserStoreAccess(ctx, dbgen.DeleteStaffUserStoreAccessParams{
		StaffUserID: targetStaffID,
		Column2:     storeIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete store access: %w", err)
	}

	// Get remaining store access for the target staff
	remainingStoreAccess, err := s.queries.GetStaffUserStoreAccess(ctx, targetStaffID)
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining store access: %w", err)
	}

	// Convert to response format
	var storeList []common.Store
	for _, access := range remainingStoreAccess {
		storeList = append(storeList, common.Store{
			ID:   utils.FormatID(access.StoreID),
			Name: access.StoreName,
		})
	}

	response := &staff.DeleteStoreAccessResponse{
		StaffUserID: utils.FormatID(targetStaffID),
		StoreList:   storeList,
	}

	return response, nil
}
