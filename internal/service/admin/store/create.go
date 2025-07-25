package adminStore

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreService struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
}

func NewCreateStoreService(queries dbgen.Querier, db *pgxpool.Pool) *CreateStoreService {
	return &CreateStoreService{
		queries: queries,
		db:      db,
	}
}

func (s *CreateStoreService) CreateStore(ctx context.Context, req adminStoreModel.CreateStoreRequest, staffContext common.StaffContext) (*adminStoreModel.CreateStoreResponse, error) {
	// Parse staff user ID
	staffUserID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.AuthStaffFailed, "invalid staff user ID", err)
	}

	// Validate role permissions (only SUPER_ADMIN and ADMIN can create stores)
	if staffContext.Role != adminStaffModel.RoleSuperAdmin && staffContext.Role != adminStaffModel.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store name already exists
	nameExists, err := s.queries.CheckStoreNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreAlreadyExists)
	}

	// Generate store ID
	storeID := utils.GenerateID()

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	// Create store
	createdStore, err := qtx.CreateStore(ctx, dbgen.CreateStoreParams{
		ID:      storeID,
		Name:    req.Name,
		Address: utils.StringPtrToPgText(&req.Address, false),
		Phone:   utils.StringPtrToPgText(&req.Phone, false),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store", err)
	}

	// If creator is ADMIN, automatically grant store access
	if staffContext.Role == adminStaffModel.RoleAdmin {
		err = qtx.CreateStaffUserStoreAccess(ctx, dbgen.CreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: staffUserID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminStoreModel.CreateStoreResponse{
		ID:       utils.FormatID(createdStore.ID),
		Name:     createdStore.Name,
		Address:  createdStore.Address.String,
		Phone:    createdStore.Phone.String,
		IsActive: createdStore.IsActive.Bool,
	}, nil
}
