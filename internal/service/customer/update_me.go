package customer

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMe struct {
	repo *sqlx.Repositories
	db   dbgen.Querier
}

func NewUpdateMe(db dbgen.Querier, repo *sqlx.Repositories) *UpdateMe {
	return &UpdateMe{
		repo: repo,
		db:   db,
	}
}

func (s *UpdateMe) UpdateMe(ctx context.Context, customerID int64, req customerModel.UpdateMeRequest) (*customerModel.UpdateMeResponse, error) {
	if !req.HasUpdates() {
		return nil, errorCodes.NewServiceError(errorCodes.ValAllFieldsEmpty, "need at least one field to update", nil)
	}

	// Validation: Parse and validate birthday format if provided
	if req.Birthday != nil {
		_, err := utils.DateStringToTime(*req.Birthday)
		if err != nil {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValFieldDateFormat)
		}
	}

	// Data Integrity: Update customer data
	result, err := s.repo.Customer.UpdateCustomer(ctx, customerID, sqlx.UpdateCustomerParams{
		Name:           req.Name,
		Phone:          req.Phone,
		Birthday:       req.Birthday,
		City:           req.City,
		Email:          req.Email,
		FavoriteShapes: req.FavoriteShapes,
		FavoriteColors: req.FavoriteColors,
		FavoriteStyles: req.FavoriteStyles,
		IsIntrovert:    req.IsIntrovert,
		CustomerNote:   req.CustomerNote,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer", err)
	}

	favoriteShapes := []string{}
	favoriteColors := []string{}
	favoriteStyles := []string{}

	if result.FavoriteShapes != nil {
		favoriteShapes = result.FavoriteShapes
	}
	if result.FavoriteColors != nil {
		favoriteColors = result.FavoriteColors
	}
	if result.FavoriteStyles != nil {
		favoriteStyles = result.FavoriteStyles
	}

	response := customerModel.UpdateMeResponse{
		ID:             utils.FormatID(result.ID),
		Name:           result.Name,
		Phone:          result.Phone,
		Birthday:       utils.PgDateToDateString(result.Birthday),
		Email:          utils.PgTextToString(result.Email),
		City:           utils.PgTextToString(result.City),
		FavoriteShapes: favoriteShapes,
		FavoriteColors: favoriteColors,
		FavoriteStyles: favoriteStyles,
		IsIntrovert:    utils.PgBoolToBool(result.IsIntrovert),
		CustomerNote:   utils.PgTextToString(result.CustomerNote),
	}

	return &response, nil
}
