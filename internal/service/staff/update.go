package staff

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStaffService struct {
	queries dbgen.Querier
	repo    sqlxRepo.StaffUserRepositoryInterface
}

func NewUpdateStaffService(queries dbgen.Querier, db *sqlx.DB) *UpdateStaffService {
	return &UpdateStaffService{
		queries: queries,
		repo:    sqlxRepo.NewStaffUserRepository(db),
	}
}

func (s *UpdateStaffService) UpdateStaff(ctx context.Context, targetID string, req staff.UpdateStaffRequest, updaterID int64, updaterRole string) (*staff.UpdateStaffResponse, error) {
	targetIDInt, err := utils.ParseID(targetID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValInputValidationFailed, "invalid target ID", err)
	}

	// Validate request has at least one field to update
	if req.Role == nil && req.IsActive == nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Validate role if provided
	if req.Role != nil && !staff.IsValidRole(*req.Role) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserInvalidRole)
	}

	// Get target staff user
	targetStaff, err := s.queries.GetStaffUserByID(ctx, targetIDInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.UserNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get target staff", err)
	}

	// Business logic validations
	if err := s.validateUpdatePermissions(updaterID, updaterRole, &targetStaff); err != nil {
		return nil, err
	}

	// Validate role change if provided
	if req.Role != nil {
		if err := s.validateRoleChange(updaterRole, *req.Role); err != nil {
			return nil, err
		}
	}

	// Perform the update
	response, err := s.repo.UpdateStaffUser(ctx, targetIDInt, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff", err)
	}

	return response, nil
}

// validateUpdatePermissions checks if the updater has permission to update the target staff
func (s *UpdateStaffService) validateUpdatePermissions(updaterID int64, updaterRole string, targetStaff *dbgen.StaffUser) error {
	// Cannot update own account
	if updaterID == targetStaff.ID {
		return errorCodes.NewServiceErrorWithCode(errorCodes.UserNotUpdateSelf)
	}

	// Cannot update SUPER_ADMIN accounts
	if targetStaff.Role == staff.RoleSuperAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Only SUPER_ADMIN and ADMIN can update staff
	if updaterRole != staff.RoleSuperAdmin && updaterRole != staff.RoleAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	return nil
}

// validateRoleChange checks if the role change is allowed
func (s *UpdateStaffService) validateRoleChange(updaterRole, newRole string) error {
	// Cannot change to SUPER_ADMIN
	if newRole == staff.RoleSuperAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	switch updaterRole {
	case staff.RoleSuperAdmin:
		// SUPER_ADMIN can change any role (except to SUPER_ADMIN, already checked above)
		return nil
	case staff.RoleAdmin:
		// ADMIN can only change MANAGER and STYLIST roles
		if newRole == staff.RoleManager || newRole == staff.RoleStylist {
			return nil
		}
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
