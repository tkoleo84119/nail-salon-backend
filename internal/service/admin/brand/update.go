package adminBrand

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBrandModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/brand"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, brandID int64, req adminBrandModel.UpdateRequest) (*adminBrandModel.UpdateResponse, error) {
	exists, err := s.queries.CheckBrandExistByID(ctx, brandID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check brand existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BrandNotFound)
	}

	if req.Name != nil {
		nameExists, err := s.queries.CheckBrandNameExistsExcludeSelf(ctx, dbgen.CheckBrandNameExistsExcludeSelfParams{
			Name: *req.Name,
			ID:   brandID,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check brand name existence", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BrandNameAlreadyExists)
		}
	}

	_, err = s.repo.Brand.UpdateBrand(ctx, brandID, sqlxRepo.UpdateBrandParams{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update brand", err)
	}

	return &adminBrandModel.UpdateResponse{
		ID: utils.FormatID(brandID),
	}, nil
}
