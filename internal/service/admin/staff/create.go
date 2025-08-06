package adminStaff

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	db      *pgxpool.Pool
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries, db *pgxpool.Pool) *Create {
	return &Create{
		db:      db,
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, req adminStaffModel.CreateParsedRequest, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateResponse, error) {
	// validate role is valid
	if !common.IsValidRole(req.Role) || req.Role == common.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffInvalidRole)
	}

	// validate permissions
	if err := s.validatePermissions(creatorRole, req.Role); err != nil {
		return nil, err
	}
	if err := s.validateStoreAccess(creatorRole, creatorStoreIDs, req.StoreIDs); err != nil {
		return nil, err
	}

	// check if username already exists
	exists, err := s.queries.CheckStaffUserExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check user existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffAlreadyExists)
	}

	// check if stores exist and active
	countInfo, err := s.queries.CheckStoresExistAndActive(ctx, req.StoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stores", err)
	}
	if countInfo.TotalCount != int64(len(req.StoreIDs)) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotFound)
	}
	if countInfo.ActiveCount != int64(len(req.StoreIDs)) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreNotActive)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to hash password", err)
	}

	// prepare of insert data
	staffID := utils.GenerateID()
	now := utils.TimeToPgTimestamptz(time.Now())

	storeAccessParams := make([]dbgen.BatchCreateStaffUserStoreAccessParams, 0, len(req.StoreIDs))
	for _, storeID := range req.StoreIDs {
		storeAccessParams = append(storeAccessParams, dbgen.BatchCreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: staffID,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	createdStaff, err := s.queries.CreateStaffUser(ctx, dbgen.CreateStaffUserParams{
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
	_, err = s.queries.BatchCreateStaffUserStoreAccess(ctx, storeAccessParams)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
	}

	_, err = s.queries.CreateStylist(ctx, dbgen.CreateStylistParams{
		ID:          staffID,
		StaffUserID: staffID,
		Name:        utils.StringPtrToPgText(&req.Username, false),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create stylist", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	response := &adminStaffModel.CreateResponse{
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
func (s *Create) validatePermissions(creatorRole, targetRole string) error {
	switch creatorRole {
	case common.RoleSuperAdmin:
		// SUPER_ADMIN can create all roles except SUPER_ADMIN
		if targetRole == common.RoleSuperAdmin {
			return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
		}
		return nil
	case common.RoleAdmin:
		// ADMIN can only create MANAGER and STYLIST
		if targetRole == common.RoleManager || targetRole == common.RoleStylist {
			return nil
		}
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	default:
		return errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}
}

// check if creator has permission to assign these stores
func (s *Create) validateStoreAccess(creatorRole string, creatorStoreIDs, targetStoreIDs []int64) error {
	// SUPER_ADMIN can assign any store
	if creatorRole == common.RoleSuperAdmin {
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
