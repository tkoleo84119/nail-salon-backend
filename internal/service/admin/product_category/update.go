package adminProductCategory

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductCategoryModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product_category"
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

func (s *Update) Update(ctx context.Context, productCategoryID int64, req adminProductCategoryModel.UpdateRequest) (*adminProductCategoryModel.UpdateResponse, error) {
	exists, err := s.queries.CheckProductCategoryExistByID(ctx, productCategoryID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product category existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CategoryNotFound)
	}

	if req.Name != nil && *req.Name != "" {
		nameExists, err := s.queries.CheckProductCategoryNameExistsExcludeSelf(ctx, dbgen.CheckProductCategoryNameExistsExcludeSelfParams{
			Name: *req.Name,
			ID:   productCategoryID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product category name existence", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CategoryNameAlreadyExists)
		}
	}

	_, err = s.repo.ProductCategory.UpdateProductCategory(ctx, productCategoryID, sqlxRepo.UpdateProductCategoryParams{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update product category", err)
	}

	return &adminProductCategoryModel.UpdateResponse{
		ID: utils.FormatID(productCategoryID),
	}, nil
}
