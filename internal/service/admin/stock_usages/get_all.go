package adminStockUsages

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewGetAll(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, req adminStockUsagesModel.GetAllParsedRequest, staffStoreIDs []int64) (*adminStockUsagesModel.GetAllResponse, error) {
	// Check store access for staff (except SUPER_ADMIN)
	if err := utils.CheckStoreAccess(storeID, staffStoreIDs); err != nil {
		return nil, err
	}

	// Get stock usages from repository
	total, results, err := s.repo.StockUsage.GetAllStockUsagesByFilter(ctx, storeID, sqlxRepo.GetAllStockUsagesByFilterParams{
		Name:    req.Name,
		IsInUse: req.IsInUse,
		Limit:   &req.Limit,
		Offset:  &req.Offset,
		Sort:    &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get stock usages", err)
	}

	items := make([]adminStockUsagesModel.GetAllListItem, len(results))
	for i, result := range results {
		items[i] = adminStockUsagesModel.GetAllListItem{
			ID: utils.FormatID(result.ID),
			Product: adminStockUsagesModel.GetAllProductItem{
				ID:   utils.FormatID(result.ProductID),
				Name: result.ProductName,
			},
			Quantity:     result.Quantity,
			IsInUse:      utils.PgBoolToBool(result.IsInUse),
			Expiration:   utils.PgDateToDateString(result.Expiration),
			UsageStarted: utils.PgDateToDateString(result.UsageStarted),
			UsageEndedAt: utils.PgDateToDateString(result.UsageEndedAt),
			CreatedAt:    utils.PgTimestamptzToTimeString(result.CreatedAt),
			UpdatedAt:    utils.PgTimestamptzToTimeString(result.UpdatedAt),
		}
	}

	return &adminStockUsagesModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
