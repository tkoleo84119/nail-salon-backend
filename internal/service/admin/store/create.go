package adminStore

import (
	"context"

	"github.com/jmoiron/sqlx"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreService struct {
	db   *sqlx.DB
	repo *sqlxRepo.Repositories
}

func NewCreateStoreService(db *sqlx.DB, repo *sqlxRepo.Repositories) *CreateStoreService {
	return &CreateStoreService{
		db:   db,
		repo: repo,
	}
}

func (s *CreateStoreService) CreateStore(ctx context.Context, req adminStoreModel.CreateStoreRequest, staffId int64, role string) (*adminStoreModel.CreateStoreResponse, error) {
	// Validate role permissions (only SUPER_ADMIN and ADMIN can create stores)
	if role != adminStaffModel.RoleSuperAdmin && role != adminStaffModel.RoleAdmin {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.AuthPermissionDenied)
	}

	// Check if store name already exists
	nameExists, err := s.repo.Store.CheckNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check store name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.StoreAlreadyExists)
	}

	// Generate store ID
	storeID := utils.GenerateID()

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Create store
	_, err = s.repo.Store.CreateStoreTx(ctx, tx, sqlxRepo.CreateStoreTxParams{
		ID:      storeID,
		Name:    req.Name,
		Address: utils.StringPtrToPgText(req.Address, true),
		Phone:   utils.StringPtrToPgText(req.Phone, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store", err)
	}

	// If creator is ADMIN, automatically grant store access
	if role == adminStaffModel.RoleAdmin {
		_, err = s.repo.StaffUserStoreAccess.CreateStaffUserStoreAccessTx(ctx, tx, sqlxRepo.CreateStaffUserStoreAccessTxParams{
			StoreID:     storeID,
			StaffUserID: staffId,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminStoreModel.CreateStoreResponse{
		ID:       utils.FormatID(storeID),
		Name:     req.Name,
		Address:  utils.PgTextToString(utils.StringPtrToPgText(req.Address, true)),
		Phone:    utils.PgTextToString(utils.StringPtrToPgText(req.Phone, true)),
		IsActive: true,
	}, nil
}
