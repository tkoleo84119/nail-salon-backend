package adminAccount

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAccountModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account"
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

func (s *Create) Create(ctx context.Context, req adminAccountModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminAccountModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(req.StoreID, creatorStoreIDs); err != nil {
		return nil, err
	}

	// Create account
	accountID := utils.GenerateID()
	err := s.queries.CreateAccount(ctx, dbgen.CreateAccountParams{
		ID:      accountID,
		StoreID: req.StoreID,
		Name:    req.Name,
		Note:    utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create account", err)
	}

	return &adminAccountModel.CreateResponse{
		ID: utils.FormatID(accountID),
	}, nil
}
