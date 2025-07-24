package customer

import (
	"context"
	"database/sql"
	"errors"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyCustomerService struct {
	queries *dbgen.Queries
}

func NewGetMyCustomerService(queries *dbgen.Queries) *GetMyCustomerService {
	return &GetMyCustomerService{
		queries: queries,
	}
}

func (s *GetMyCustomerService) GetMyCustomer(ctx context.Context, customerContext common.CustomerContext) (*customer.GetMyCustomerResponse, error) {
	// Use customer ID directly from context
	customerID := customerContext.CustomerID

	// Get customer data
	customerData, err := s.queries.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer", err)
	}

	// Build response
	response := &customer.GetMyCustomerResponse{
		ID:           utils.FormatID(customerData.ID),
		Name:         customerData.Name,
		Phone:        customerData.Phone,
		Birthday:     utils.PgDateToDateString(customerData.Birthday),
		City:         utils.PgTextToString(customerData.City),
		IsIntrovert:  customerData.IsIntrovert.Bool,
		Referrer:     utils.PgTextToString(customerData.Referrer),
		CustomerNote: utils.PgTextToString(customerData.CustomerNote),
	}

	// Handle optional fields
	if len(customerData.FavoriteShapes) > 0 {
		response.FavoriteShapes = &customerData.FavoriteShapes
	}

	if len(customerData.FavoriteColors) > 0 {
		response.FavoriteColors = &customerData.FavoriteColors
	}

	if len(customerData.FavoriteStyles) > 0 {
		response.FavoriteStyles = &customerData.FavoriteStyles
	}

	if len(customerData.ReferralSource) > 0 {
		response.ReferralSource = &customerData.ReferralSource
	}

	return response, nil
}
