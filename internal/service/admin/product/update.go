package adminProduct

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, storeID, productID int64, req adminProductModel.UpdateParsedRequest, creatorStoreIDs []int64) (*adminProductModel.UpdateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// check product exists and belongs to the store
	product, err := s.queries.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product existence", err)
	}
	if product.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
	}

	// check brand exists
	if req.BrandID != nil {
		brandExists, err := s.queries.CheckBrandExistByID(ctx, *req.BrandID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check brand existence", err)
		}
		if !brandExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BrandNotFound)
		}
	}

	// check category exists
	if req.CategoryID != nil {
		categoryExists, err := s.queries.CheckProductCategoryExistByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check category existence", err)
		}
		if !categoryExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CategoryNotFound)
		}
	}

	// check name and brand uniqueness
	if (req.Name != nil && *req.Name != "") || req.BrandID != nil {
		// decide which name and brandID to check
		checkName := product.Name
		if req.Name != nil && *req.Name != "" {
			checkName = *req.Name
		}

		checkBrandID := product.BrandID
		if req.BrandID != nil {
			checkBrandID = *req.BrandID
		}

		// check uniqueness
		nameExists, err := s.queries.CheckProductNameBrandExistsInStoreExcluding(ctx, dbgen.CheckProductNameBrandExistsInStoreExcludingParams{
			ID:      productID,
			StoreID: storeID,
			Name:    checkName,
			BrandID: checkBrandID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product name uniqueness", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNameBrandAlreadyExistsInStore)
		}
	}

	updateResult, err := s.repo.Product.UpdateStoreProduct(ctx, productID, sqlxRepo.UpdateStoreProductParams{
		BrandID:         req.BrandID,
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		CurrentStock:    req.CurrentStock,
		SafetyStock:     req.SafetyStock,
		Unit:            req.Unit,
		StorageLocation: req.StorageLocation,
		Note:            req.Note,
		IsActive:        req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update product", err)
	}

	return &adminProductModel.UpdateResponse{
		ID: utils.FormatID(updateResult.ID),
	}, nil
}
