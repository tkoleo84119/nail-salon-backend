package staff

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
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

func (s *CreateStaffService) CreateStaff(ctx context.Context, req staff.CreateStaffRequest, creatorRole string, creatorStoreIDs []int64) (*staff.CreateStaffResponse, error) {
	if err := s.validatePermissions(creatorRole, req.Role); err != nil {
		return nil, err
	}

	if !staff.IsValidRole(req.Role) {
		return nil, fmt.Errorf("invalid role: %s", req.Role)
	}

	if req.Role == staff.RoleSuperAdmin {
		return nil, fmt.Errorf("cannot create SUPER_ADMIN role")
	}

	if err := s.validateStoreAccess(creatorRole, creatorStoreIDs, req.StoreIDs); err != nil {
		return nil, err
	}

	// check if username or email already exists
	exists, err := s.queries.CheckStaffUserExists(ctx, dbgen.CheckStaffUserExistsParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username or email already exists")
	}

	// check if stores exist and active
	storeCheck, err := s.queries.CheckStoresExistAndActive(ctx, req.StoreIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to check stores: %w", err)
	}
	if storeCheck.TotalCount != int64(len(req.StoreIDs)) {
		return nil, fmt.Errorf("some stores do not exist")
	}
	if storeCheck.ActiveCount != storeCheck.TotalCount {
		return nil, fmt.Errorf("some stores are not active")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := dbgen.New(tx)

	staffID := utils.GenerateID()
	createdStaff, err := qtx.CreateStaffUser(ctx, dbgen.CreateStaffUserParams{
		ID:           staffID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create staff user: %w", err)
	}

	// batch create store access records
	err = qtx.BatchCreateStaffUserStoreAccess(ctx, dbgen.BatchCreateStaffUserStoreAccessParams{
		Column1:     req.StoreIDs,
		StaffUserID: staffID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create store access: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// get store list information
	stores, err := s.queries.GetStoresByIDs(ctx, req.StoreIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get store information: %w", err)
	}

	storeList := make([]common.Store, 0, len(stores))
	for _, store := range stores {
		storeList = append(storeList, common.Store{
			ID:   store.ID,
			Name: store.Name,
		})
	}

	response := &staff.CreateStaffResponse{
		ID:        strconv.FormatInt(createdStaff.ID, 10),
		Username:  createdStaff.Username,
		Email:     createdStaff.Email,
		Role:      createdStaff.Role,
		StoreList: storeList,
	}

	return response, nil
}

// check if creator has permission to create staff with target role
func (s *CreateStaffService) validatePermissions(creatorRole, targetRole string) error {
	switch creatorRole {
	case staff.RoleSuperAdmin:
		// SUPER_ADMIN can create all roles except SUPER_ADMIN
		if targetRole == staff.RoleSuperAdmin {
			return fmt.Errorf("SUPER_ADMIN cannot create another SUPER_ADMIN")
		}
		return nil
	case staff.RoleAdmin:
		// ADMIN can only create MANAGER and STYLIST
		if targetRole == staff.RoleManager || targetRole == staff.RoleStylist {
			return nil
		}
		return fmt.Errorf("ADMIN can only create MANAGER or STYLIST roles")
	default:
		return fmt.Errorf("insufficient permissions to create staff")
	}
}

// check if creator has permission to assign these stores
func (s *CreateStaffService) validateStoreAccess(creatorRole string, creatorStoreIDs, targetStoreIDs []int64) error {
	// SUPER_ADMIN can assign any store
	if creatorRole == staff.RoleSuperAdmin {
		return nil
	}

	// for ADMIN, check if has permission to assign these stores
	creatorStoreMap := make(map[int64]bool)
	for _, storeID := range creatorStoreIDs {
		creatorStoreMap[storeID] = true
	}

	for _, targetStoreID := range targetStoreIDs {
		if !creatorStoreMap[targetStoreID] {
			return fmt.Errorf("no permission to assign store ID: %d", targetStoreID)
		}
	}

	return nil
}