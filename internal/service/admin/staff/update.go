package adminStaff

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries   *dbgen.Queries
	repo      *sqlxRepo.Repositories
	authCache cache.AuthCacheInterface
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, authCache cache.AuthCacheInterface) *Update {
	return &Update{
		queries:   queries,
		repo:      repo,
		authCache: authCache,
	}
}

func (s *Update) Update(ctx context.Context, staffID int64, req adminStaffModel.UpdateRequest, updaterID int64, updaterRole string) (*adminStaffModel.UpdateResponse, error) {
	// Validate role if provided
	if req.Role != nil && !common.IsValidRole(*req.Role) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffInvalidRole)
	}

	// Get target staff user
	targetStaff, err := s.queries.GetStaffUserByID(ctx, staffID)
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
	updatedStaff, err := s.repo.Staff.UpdateStaffUser(ctx, staffID, sqlxRepo.UpdateStaffUserParams{
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff user", err)
	}

	// delete staff context from cache
	if cacheErr := s.authCache.DeleteStaffContext(ctx, staffID); cacheErr != nil {
		log.Println("failed to delete staff context from cache", cacheErr)
	}

	response := &adminStaffModel.UpdateResponse{
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
func (s *Update) validateUpdatePermissions(updaterID int64, updaterRole string, targetStaff dbgen.StaffUser) error {
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
func (s *Update) validateRoleChange(updaterRole, newRole string) error {
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
