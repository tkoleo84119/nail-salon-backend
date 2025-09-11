package adminProductCategory

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductCategoryModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product_category"
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

func (s *Create) Create(ctx context.Context, req adminProductCategoryModel.CreateRequest) (*adminProductCategoryModel.CreateResponse, error) {
	nameExists, err := s.queries.CheckProductCategoryNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product category name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CategoryNameAlreadyExists)
	}

	categoryID := utils.GenerateID()
	_, err = s.queries.CreateProductCategory(ctx, dbgen.CreateProductCategoryParams{
		ID:   categoryID,
		Name: req.Name,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create product category", err)
	}

	return &adminProductCategoryModel.CreateResponse{
		ID: utils.FormatID(categoryID),
	}, nil
}
