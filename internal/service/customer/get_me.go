package customer

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMe struct {
	queries *dbgen.Queries
}

func NewGetMe(queries *dbgen.Queries) GetMeInterface {
	return &GetMe{
		queries: queries,
	}
}

func (s *GetMe) GetMe(ctx context.Context, customerID int64) (*customerModel.GetMeResponse, error) {
	// Get customer data
	customerData, err := s.queries.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer", err)
	}

	favoriteShapes := []string{}
	favoriteColors := []string{}
	favoriteStyles := []string{}
	if customerData.FavoriteShapes != nil {
		favoriteShapes = customerData.FavoriteShapes
	}
	if customerData.FavoriteColors != nil {
		favoriteColors = customerData.FavoriteColors
	}
	if customerData.FavoriteStyles != nil {
		favoriteStyles = customerData.FavoriteStyles
	}

	// Build response
	response := &customerModel.GetMeResponse{
		ID:             utils.FormatID(customerData.ID),
		Name:           customerData.Name,
		Phone:          customerData.Phone,
		Email:          utils.PgTextToString(customerData.Email),
		Birthday:       utils.PgDateToDateString(customerData.Birthday),
		City:           utils.PgTextToString(customerData.City),
		FavoriteShapes: favoriteShapes,
		FavoriteColors: favoriteColors,
		FavoriteStyles: favoriteStyles,
		IsIntrovert:    customerData.IsIntrovert.Bool,
		CustomerNote:   utils.PgTextToString(customerData.CustomerNote),
	}

	return response, nil
}
