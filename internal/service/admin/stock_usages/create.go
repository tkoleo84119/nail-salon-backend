package adminStockUsages

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
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

func (s *Create) Create(ctx context.Context, storeID int64, req adminStockUsagesModel.CreateParsedRequest, creatorStoreIDs []int64) (*adminStockUsagesModel.CreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, creatorStoreIDs); err != nil {
		return nil, err
	}

	product, err := s.queries.GetProductByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check product existence", err)
	}
	if product.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotBelongToStore)
	}

	// check if stock usage more than product stock
	if product.CurrentStock < int32(req.Quantity) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductStockNotEnough)
	}

	// Create stock usage
	stockUsageID := utils.GenerateID()
	var expiration pgtype.Date
	if req.Expiration != nil {
		expiration = utils.TimeToPgDate(*req.Expiration)
	} else {
		expiration = pgtype.Date{Valid: false}
	}

	err = s.queries.CreateStockUsage(ctx, dbgen.CreateStockUsageParams{
		ID:           stockUsageID,
		ProductID:    req.ProductID,
		Quantity:     int32(req.Quantity),
		Expiration:   expiration,
		UsageStarted: utils.TimeToPgDate(req.UsageStarted),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create stock usage", err)
	}

	return &adminStockUsagesModel.CreateResponse{
		ID: utils.FormatID(stockUsageID),
	}, nil
}
