package adminSupplier

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
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

func (s *Create) Create(ctx context.Context, req adminSupplierModel.CreateRequest) (*adminSupplierModel.CreateResponse, error) {
	nameExists, err := s.queries.CheckSupplierNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check supplier name existence", err)
	}
	if nameExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SupplierNameAlreadyExists)
	}

	supplierID := utils.GenerateID()
	_, err = s.queries.CreateSupplier(ctx, dbgen.CreateSupplierParams{
		ID:   supplierID,
		Name: req.Name,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create supplier", err)
	}

	return &adminSupplierModel.CreateResponse{
		ID: utils.FormatID(supplierID),
	}, nil
}
