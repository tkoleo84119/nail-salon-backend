package adminCustomer

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries dbgen.Querier
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries dbgen.Querier, repo *sqlxRepo.Repositories) *Update {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, customerID int64, req adminCustomerModel.UpdateRequest) (*adminCustomerModel.UpdateResponse, error) {
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValAllFieldsEmpty)
	}

	// verify customer exists
	exists, err := s.queries.CheckCustomerExistsByID(ctx, customerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer exists", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
	}

	// update customer
	customer, err := s.repo.Customer.UpdateCustomer(ctx, customerID, sqlxRepo.UpdateCustomerParams{
		StoreNote:     req.StoreNote,
		Level:         req.Level,
		IsBlacklisted: req.IsBlacklisted,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer", err)
	}

	return &adminCustomerModel.UpdateResponse{
		ID:            utils.FormatID(customer.ID),
		Name:          customer.Name,
		Phone:         customer.Phone,
		Birthday:      utils.PgDateToDateString(customer.Birthday),
		City:          utils.PgTextToString(customer.City),
		Level:         utils.PgTextToString(customer.Level),
		IsBlacklisted: utils.PgBoolToBool(customer.IsBlacklisted),
		LastVisitAt:   utils.PgTimestamptzToTimeString(customer.LastVisitAt),
		UpdatedAt:     utils.PgTimestamptzToTimeString(customer.UpdatedAt),
	}, nil
}
