package adminStaff

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
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
func (s *CreateStoreAccessService) CreateStoreAccess(ctx context.Context, targetID string, req adminStaffModel.CreateStoreAccessRequest, creatorID int64, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStoreAccessResponse, bool, error) {
	// Parse target staff ID
	targetStaffID, err := utils.ParseID(targetID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid target staff ID", err)
	}

	// Parse store ID
	storeID, err := utils.ParseID(req.StoreID)
	if err != nil {
		return nil, false, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid store ID", err)
	}

	// Check if target staff exists
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetStaffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, false, fmt.Errorf("failed to get target staff: %w", err)
	}

	// Cannot modify self
	if targetStaffID == creatorID {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotUpdateSelf)
	}

	// Cannot modify SUPER_ADMIN
	if targetStaff.Role == adminStaffModel.RoleSuperAdmin {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store exists and is active
	store, err := s.queries.GetStoreByID(ctx, storeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffStoreNotFound)
		}
		return nil, false, fmt.Errorf("failed to get store: %w", err)
	}
	if !store.IsActive.Bool {
		return nil, false, errorCodes.NewServiceErrorWithCode(errorCodes.StaffStoreNotActive)
	}

	// Check if creator has access to this store (except SUPER_ADMIN)
	if creatorRole != adminStaffModel.RoleSuperAdmin {
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

	response := &adminStaffModel.CreateStoreAccessResponse{
		StaffUserID: targetID,
		StoreList:   storeList,
	}

	return response, isNewlyCreated, nil
}
