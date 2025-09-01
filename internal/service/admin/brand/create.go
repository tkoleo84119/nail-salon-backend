package adminBrand

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBrandModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/brand"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, req adminBrandModel.CreateRequest) (*adminBrandModel.CreateResponse, error) {
	nameExists, err := s.queries.CheckBrandNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check brand name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BrandNameAlreadyExists)
	}

	brandID := utils.GenerateID()
	brand, err := s.queries.CreateBrand(ctx, dbgen.CreateBrandParams{
		ID:   brandID,
		Name: req.Name,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create brand", err)
	}

	return &adminBrandModel.CreateResponse{
		ID: utils.FormatID(brand.ID),
	}, nil
}
