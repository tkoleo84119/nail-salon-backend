package customer

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyCustomerService struct {
	customerRepo sqlx.CustomerRepositoryInterface
	db           dbgen.Querier
}

func NewUpdateMyCustomerService(db dbgen.Querier, customerRepo sqlx.CustomerRepositoryInterface) *UpdateMyCustomerService {
	return &UpdateMyCustomerService{
		customerRepo: customerRepo,
		db:           db,
	}
}

func (s *UpdateMyCustomerService) UpdateMyCustomer(ctx context.Context, customerID int64, req customerModel.UpdateMyCustomerRequest) (*customerModel.UpdateMyCustomerResponse, error) {
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	// Validation: Parse and validate birthday format if provided
	if req.Birthday != nil {
		_, err := utils.DateStringToTime(*req.Birthday)
		if err != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValDateFormatInvalid)
		}
	}

	// Check if customer exists
	_, err := s.db.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
	}

	// Data Integrity: Update customer data
	result, err := s.customerRepo.UpdateMyCustomer(ctx, customerID, req)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer", err)
	}

	return result, nil
}
