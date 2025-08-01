package adminStaff

import (
	"context"

	"github.com/jmoiron/sqlx"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStaffService struct {
	db   *sqlx.DB
	repo *sqlxRepo.Repositories
}

func NewCreateStaffService(db *sqlx.DB, repo *sqlxRepo.Repositories) *CreateStaffService {
	return &CreateStaffService{
		db:   db,
		repo: repo,
	}
}

func (s *CreateStaffService) CreateStaff(ctx context.Context, req adminStaffModel.CreateStaffParsedRequest, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStaffResponse, error) {
	// validate role is valid
	if !common.IsValidRole(req.Role) || req.Role == adminStaffModel.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffInvalidRole)
	}

	// validate permissions
	if err := s.validatePermissions(creatorRole, req.Role); err != nil {
		return nil, err
	}
	if err := s.validateStoreAccess(creatorRole, creatorStoreIDs, req.StoreIDs); err != nil {
		return nil, err
	}

	// check if username or email already exists
	exists, err := s.repo.Staff.CheckStaffUserExists(ctx, req.Username)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check user existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffAlreadyExists)
	}

	// check if stores exist and active
	storeCount, err := s.repo.Store.CheckStoresExistAndActive(ctx, req.StoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stores", err)
	}
	if storeCount != len(req.StoreIDs) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to hash password", err)
	}

	// prepare of insert data
	staffID := utils.GenerateID()

	storeAccessParams := make([]sqlxRepo.CreateStaffUserStoreAccessTxParams, 0, len(req.StoreIDs))
	for _, storeID := range req.StoreIDs {
		storeAccessParams = append(storeAccessParams, sqlxRepo.CreateStaffUserStoreAccessTxParams{
			StoreID:     storeID,
			StaffUserID: staffID,
		})
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	createdStaff, err := s.repo.Staff.CreateStaffUserTx(ctx, tx, sqlxRepo.CreateStaffUserTxParams{
		ID:           staffID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create staff user", err)
	}

	// batch create store access records
	err = s.repo.StaffUserStoreAccess.BatchCreateStaffUserStoreAccessTx(ctx, tx, storeAccessParams)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
	}

	_, err = s.repo.Stylist.CreateStylistTx(ctx, tx, sqlxRepo.CreateStylistTxParams{
		ID:          staffID,
		StaffUserID: staffID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create stylist", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	response := &adminStaffModel.CreateStaffResponse{
		ID:        utils.FormatID(createdStaff.ID),
		Username:  createdStaff.Username,
		Email:     createdStaff.Email,
		Role:      createdStaff.Role,
		IsActive:  true,
		CreatedAt: utils.PgTimestamptzToTimeString(createdStaff.CreatedAt),
		UpdatedAt: utils.PgTimestamptzToTimeString(createdStaff.UpdatedAt),
	}

	return response, nil
}

// check if creator has permission to create staff with target role
func (s *CreateStaffService) validatePermissions(creatorRole, targetRole string) error {
	switch creatorRole {
	case adminStaffModel.RoleSuperAdmin:
		// SUPER_ADMIN can create all roles except SUPER_ADMIN
		if targetRole == adminStaffModel.RoleSuperAdmin {
			return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
		return nil
	case adminStaffModel.RoleAdmin:
		// ADMIN can only create MANAGER and STYLIST
		if targetRole == adminStaffModel.RoleManager || targetRole == adminStaffModel.RoleStylist {
			return nil
		}
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}

// check if creator has permission to assign these stores
func (s *CreateStaffService) validateStoreAccess(creatorRole string, creatorStoreIDs, targetStoreIDs []int64) error {
	// SUPER_ADMIN can assign any store
	if creatorRole == adminStaffModel.RoleSuperAdmin {
		return nil
	}

	// for ADMIN, check if has permission to assign these stores
	creatorStoreMap := make(map[int64]bool)
	for _, storeID := range creatorStoreIDs {
		creatorStoreMap[storeID] = true
	}

	for _, targetStoreID := range targetStoreIDs {
		if !creatorStoreMap[targetStoreID] {
			return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
	}

	return nil
}
