package adminStore

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	db        *pgxpool.Pool
	queries   *dbgen.Queries
	authCache cache.AuthCacheInterface
}

func NewCreate(queries *dbgen.Queries, db *pgxpool.Pool, authCache cache.AuthCacheInterface) CreateInterface {
	return &Create{
		db:        db,
		queries:   queries,
		authCache: authCache,
	}
}

func (s *Create) Create(ctx context.Context, req adminStoreModel.CreateRequest, staffId int64, role string) (*adminStoreModel.CreateResponse, error) {
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

	qtx := s.queries.WithTx(tx)

	// Create store
	err = qtx.CreateStore(ctx, dbgen.CreateStoreParams{
		ID:      storeID,
		Name:    req.Name,
		Address: utils.StringPtrToPgText(req.Address, true),
		Phone:   utils.StringPtrToPgText(req.Phone, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store", err)
	}

	// If creator is ADMIN, automatically grant store access
	if role == common.RoleAdmin {
		err = qtx.CreateStaffUserStoreAccess(ctx, dbgen.CreateStaffUserStoreAccessParams{
			StoreID:     storeID,
			StaffUserID: staffId,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create store access", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// delete all staff context from cache
	if cacheErr := s.authCache.DeleteAllStaffContext(ctx); cacheErr != nil {
		log.Println("failed to delete staff context from cache", cacheErr)
	}

	return &adminStoreModel.CreateResponse{
		ID: utils.FormatID(storeID),
	}, nil
}
