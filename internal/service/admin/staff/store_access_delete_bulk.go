package adminStaff

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteStoreAccessBulkService struct {
	queries dbgen.Querier
}

func NewDeleteStoreAccessBulkService(queries dbgen.Querier) *DeleteStoreAccessBulkService {
	return &DeleteStoreAccessBulkService{
		queries: queries,
	}
}

// DeleteStoreAccessBulk deletes store access for a staff member
func (s *DeleteStoreAccessBulkService) DeleteStoreAccessBulk(ctx context.Context, targetID string, req adminStaffModel.DeleteStoreAccessBulkRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.DeleteStoreAccessBulkResponse, error) {
	// Parse target staff ID
	targetStaffID, err := utils.ParseID(targetID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid target staff ID", err)
	}

	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetStaffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserNotFound)
		}
		return nil, fmt.Errorf("failed to get target staff: %w", err)
	}

	// Cannot modify self
	if targetStaffID == creatorID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserNotUpdateSelf)
	}
	// Cannot modify SUPER_ADMIN store access
	if targetStaff.Role == adminStaffModel.RoleSuperAdmin {
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
	if creatorRole != adminStaffModel.RoleSuperAdmin {
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

	response := &adminStaffModel.DeleteStoreAccessBulkResponse{
		StaffUserID: utils.FormatID(targetStaffID),
		StoreList:   storeList,
	}

	return response, nil
}
