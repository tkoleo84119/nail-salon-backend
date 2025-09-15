package adminProduct

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) CreateInterface {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, storeID int64, req adminProductModel.CreateParsedRequest, role string, creatorStoreIDs []int64) (*adminProductModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs, role); err != nil {
		return nil, err
	}

	brandExists, err := s.queries.CheckBrandExistByID(ctx, req.BrandID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check brand existence", err)
	}
	if !brandExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BrandNotFound)
	}

	categoryExists, err := s.queries.CheckProductCategoryExistByID(ctx, req.CategoryID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check category existence", err)
	}
	if !categoryExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CategoryNotFound)
	}

	// check if product name and brand already exists in one store
	nameExists, err := s.queries.CheckProductNameBrandExistsInStore(ctx, dbgen.CheckProductNameBrandExistsInStoreParams{
		StoreID: storeID,
		Name:    req.Name,
		BrandID: req.BrandID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNameBrandAlreadyExistsInStore)
	}

	productID := utils.GenerateID()
	err = s.queries.CreateProduct(ctx, dbgen.CreateProductParams{
		ID:              productID,
		StoreID:         storeID,
		Name:            req.Name,
		BrandID:         req.BrandID,
		CategoryID:      req.CategoryID,
		CurrentStock:    req.CurrentStock,
		SafetyStock:     utils.Int32PtrToPgInt4(&req.SafetyStock),
		Unit:            utils.StringPtrToPgText(req.Unit, true),
		StorageLocation: utils.StringPtrToPgText(req.StorageLocation, true),
		Note:            utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create product", err)
	}

	return &adminProductModel.CreateResponse{
		ID: utils.FormatID(productID),
	}, nil
}
