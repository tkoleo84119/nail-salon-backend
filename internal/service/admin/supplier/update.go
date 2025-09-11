package adminSupplier

import (
	"context"
	"strings"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
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

func (s *Update) Update(ctx context.Context, id int64, req adminSupplierModel.UpdateRequest) (*adminSupplierModel.UpdateResponse, error) {
	exists, err := s.queries.CheckSupplierExistsByID(ctx, id)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check supplier existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNotFound)
	}

	// If updating name, check if the name is unique (excluding self)
	if req.Name != nil && *req.Name != "" {
		*req.Name = strings.TrimSpace(*req.Name)

		nameExists, err := s.queries.CheckSupplierNameExistsExcluding(ctx, dbgen.CheckSupplierNameExistsExcludingParams{
			Name: *req.Name,
			ID:   id,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to check supplier name uniqueness", err)
		}
		if nameExists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNameAlreadyExists)
		}
	}

	// Update supplier
	result, err := s.repo.Supplier.UpdateSupplier(ctx, id, sqlxRepo.UpdateSupplierParams{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update supplier", err)
	}

	response := &adminSupplierModel.UpdateResponse{
		ID: utils.FormatID(result.ID),
	}

	return response, nil
}
