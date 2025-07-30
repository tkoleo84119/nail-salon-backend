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

type CreateStaffService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreateStaffService(queries dbgen.Querier, db *pgxpool.Pool) *CreateStaffService {
	return &CreateStaffService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateStaffService) CreateStaff(ctx context.Context, req adminStaffModel.CreateStaffRequest, creatorRole string, creatorStoreIDs []int64) (*adminStaffModel.CreateStaffResponse, error) {
	// Convert string store IDs to int64
	storeIDs, err := utils.ParseIDSlice(req.StoreIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "invalid store IDs", err)
	}

	// validate role is valid
	if !adminStaffModel.IsValidRole(req.Role) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffInvalidRole)
	}
	if req.Role == adminStaffModel.RoleSuperAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// validate permissions
	if err := s.validatePermissions(creatorRole, req.Role); err != nil {
		return nil, err
	}
	if err := s.validateStoreAccess(creatorRole, creatorStoreIDs, storeIDs); err != nil {
		return nil, err
	}

	// check if username or email already exists
	exists, err := s.queries.CheckStaffUserExists(ctx, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check user existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffAlreadyExists)
	}

	// check if stores exist and active
	storeCheck, err := s.queries.CheckStoresExistAndActive(ctx, storeIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check stores", err)
	}
	if storeCheck.TotalCount != int64(len(storeIDs)) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffStoreNotFound)
	}
	if storeCheck.ActiveCount != storeCheck.TotalCount {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StaffStoreNotActive)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to hash password", err)
	}

	// prepare of insert data
	now := time.Now()
	staffID := utils.GenerateID()

	storeAccessParams := make([]dbgen.BatchCreateStaffUserStoreAccessParams, 0, len(storeIDs))
	for _, storeID := range storeIDs {
		storeAccessParams = append(storeAccessParams, dbgen.BatchCreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: staffID,
			CreatedAt:   utils.TimeToPgTimestamptz(now),
			UpdatedAt:   utils.TimeToPgTimestamptz(now),
		})
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	createdStaff, err := qtx.CreateStaffUser(ctx, dbgen.CreateStaffUserParams{
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
	_, err = qtx.BatchCreateStaffUserStoreAccess(ctx, storeAccessParams)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// get store list information
	stores, err := s.queries.GetStoresByIDs(ctx, storeIDs)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get store information", err)
	}

	storeList := make([]common.Store, 0, len(stores))
	for _, store := range stores {
		storeList = append(storeList, common.Store{
			ID:   utils.FormatID(store.ID),
			Name: store.Name,
		})
	}

	response := &adminStaffModel.CreateStaffResponse{
		ID:        utils.FormatID(createdStaff.ID),
		Username:  createdStaff.Username,
		Email:     createdStaff.Email,
		Role:      createdStaff.Role,
		IsActive:  createdStaff.IsActive.Bool,
		StoreList: storeList,
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
