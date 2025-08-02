package adminStaff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/jmoiron/sqlx"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStaffService struct {
	repo sqlxRepo.StaffUserRepositoryInterface
}

func NewUpdateStaffService(db *sqlx.DB) *UpdateStaffService {
	return &UpdateStaffService{
		repo: sqlxRepo.NewStaffUserRepository(db),
	}
}

func (s *UpdateStaffService) UpdateStaff(ctx context.Context, staffID int64, req adminStaffModel.UpdateStaffRequest, updaterID int64, updaterRole string) (*adminStaffModel.UpdateStaffResponse, error) {
	// validate request has at least one field to update
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// Validate role if provided
	if req.Role != nil && !common.IsValidRole(*req.Role) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffInvalidRole)
	}

	// Get target staff user
	targetStaff, err := s.repo.GetStaffUserByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get target staff", err)
	}

	// Business logic validations
	if err := s.validateUpdatePermissions(updaterID, updaterRole, targetStaff); err != nil {
		return nil, err
	}

	// Validate role change if provided
	if req.Role != nil {
		if err := s.validateRoleChange(updaterRole, *req.Role); err != nil {
			return nil, err
		}
	}

	// Perform the update
	updatedStaff, err := s.repo.UpdateStaffUser(ctx, staffID, sqlxRepo.UpdateStaffUserParams{
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff user", err)
	}

	response := &adminStaffModel.UpdateStaffResponse{
		ID:        utils.FormatID(updatedStaff.ID),
		Username:  updatedStaff.Username,
		Email:     updatedStaff.Email,
		Role:      updatedStaff.Role,
		IsActive:  utils.PgBoolToBool(updatedStaff.IsActive),
		CreatedAt: utils.PgTimestamptzToTimeString(updatedStaff.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(updatedStaff.UpdatedAt),
	}

	return response, nil
}

// validateUpdatePermissions checks if the updater has permission to update the target staff
func (s *UpdateStaffService) validateUpdatePermissions(updaterID int64, updaterRole string, targetStaff *sqlxRepo.GetStaffUserByIDResponse) error {
	// Cannot update own account
	if updaterID == targetStaff.ID {
		return errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotUpdateSelf)
	}

	// Cannot update SUPER_ADMIN accounts
	if targetStaff.Role == common.RoleSuperAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Only SUPER_ADMIN and ADMIN can update staff
	if updaterRole != common.RoleSuperAdmin && updaterRole != common.RoleAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	return nil
}

// validateRoleChange checks if the role change is allowed
func (s *UpdateStaffService) validateRoleChange(updaterRole, newRole string) error {
	// Cannot change to SUPER_ADMIN
	if newRole == common.RoleSuperAdmin {
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	switch updaterRole {
	case common.RoleSuperAdmin:
		// SUPER_ADMIN can change any role (except to SUPER_ADMIN, already checked above)
		return nil
	case common.RoleAdmin:
		// ADMIN can only change MANAGER and STYLIST roles
		if newRole == common.RoleManager || newRole == common.RoleStylist {
			return nil
		}
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}
