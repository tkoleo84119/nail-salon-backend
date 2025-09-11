package adminCustomer

import (
	"context"
	"log"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries   *dbgen.Queries
	repo      *sqlxRepo.Repositories
	authCache cache.AuthCacheInterface
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories, authCache cache.AuthCacheInterface) UpdateInterface {
	return &Update{
		queries:   queries,
		repo:      repo,
		authCache: authCache,
	}
}

func (s *Update) Update(ctx context.Context, customerID int64, req adminCustomerModel.UpdateRequest) (*adminCustomerModel.UpdateResponse, error) {
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

	if cacheErr := s.authCache.DeleteCustomerContext(ctx, customerID); cacheErr != nil {
		log.Println("failed to delete customer context from cache", cacheErr)
	}

	return &adminCustomerModel.UpdateResponse{
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
	}, nil
}
