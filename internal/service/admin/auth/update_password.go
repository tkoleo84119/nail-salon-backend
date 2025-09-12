package adminAuth

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	cacheService "github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdatePassword struct {
	queries dbgen.Querier
	cache   cacheService.AuthCacheInterface
}

// NewUpdatePassword creates a new update password service
func NewUpdatePassword(queries dbgen.Querier, cache cacheService.AuthCacheInterface) UpdatePasswordInterface {
	return &UpdatePassword{
		queries: queries,
		cache:   cache,
	}
}

// UpdatePassword updates staff user password
func (s *UpdatePassword) UpdatePassword(ctx context.Context, req adminAuthModel.UpdatePasswordParsedRequest, staffContext *common.StaffContext) (*adminAuthModel.UpdatePasswordResponse, error) {
	// Permission check: SUPER_ADMIN can update any password, others can only update their own
	if staffContext.Role != common.RoleSuperAdmin && staffContext.UserID != req.StaffId {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// if not SUPER_ADMIN, need to pass oldPassword
	if staffContext.Role != common.RoleSuperAdmin && req.OldPassword == nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
	}

	// Get target staff user by ID
	targetStaff, err := s.queries.GetStaffUserByID(ctx, req.StaffId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffNotFound)
		}
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SysDatabaseError)
	}

	// Verify old password
	if req.OldPassword != nil {
		if !utils.CheckPassword(*req.OldPassword, targetStaff.PasswordHash) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthInvalidCredentials)
		}
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to hash password", err)
	}

	// Update password in database
	updatedStaffID, err := s.queries.UpdateStaffUserPassword(ctx, dbgen.UpdateStaffUserPasswordParams{
		ID:           req.StaffId,
		PasswordHash: newPasswordHash,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update staff user password", err)
	}

	// Clear staff context cache to force re-authentication
	if cacheErr := s.cache.DeleteStaffContext(ctx, req.StaffId); cacheErr != nil {
		log.Println("failed to delete staff context from cache", cacheErr)
	}

	return &adminAuthModel.UpdatePasswordResponse{
		ID: utils.FormatID(updatedStaffID),
	}, nil
}
