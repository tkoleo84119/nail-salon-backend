package adminCustomer

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, customerID int64) (*adminCustomerModel.GetResponse, error) {
	customer, err := s.queries.GetCustomerByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
		}
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.SysDatabaseError)
	}

	// Convert to response format
	response := &adminCustomerModel.GetResponse{
		ID:             utils.FormatID(customer.ID),
		Name:           customer.Name,
		LineName:       utils.PgTextToString(customer.LineName),
		Phone:          customer.Phone,
		Birthday:       utils.PgDateToDateString(customer.Birthday),
		Email:          utils.PgTextToString(customer.Email),
		City:           utils.PgTextToString(customer.City),
		FavoriteShapes: customer.FavoriteShapes,
		FavoriteColors: customer.FavoriteColors,
		FavoriteStyles: customer.FavoriteStyles,
		IsIntrovert:    utils.PgBoolToBool(customer.IsIntrovert),
		ReferralSource: customer.ReferralSource,
		Referrer:       utils.PgTextToString(customer.Referrer),
		CustomerNote:   utils.PgTextToString(customer.CustomerNote),
		StoreNote:      utils.PgTextToString(customer.StoreNote),
		Level:          utils.PgTextToString(customer.Level),
		IsBlacklisted:  utils.PgBoolToBool(customer.IsBlacklisted),
		LastVisitAt:    utils.PgTimestamptzToTimeString(customer.LastVisitAt),
		CreatedAt:      utils.PgTimestamptzToTimeString(customer.CreatedAt),
		UpdatedAt:      utils.PgTimestamptzToTimeString(customer.UpdatedAt),
	}

	return response, nil
}
