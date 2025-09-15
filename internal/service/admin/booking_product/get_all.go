package adminBookingProduct

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking_product"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewGetAll(queries *dbgen.Queries, repo *sqlxRepo.Repositories) GetAllInterface {
	return &GetAll{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.GetAllParsedRequest, role string, staffStoreIDs []int64) (*adminBookingProductModel.GetAllResponse, error) {
	// Check store access permission
	if err := utils.CheckStoreAccess(storeID, staffStoreIDs, role); err != nil {
		return nil, err
	}

	// Verify booking exists
	booking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking", err)
	}

	// Verify booking belongs to the store
	if booking.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
	}

	// Query booking products with filtering and pagination
	total, items, err := s.repo.BookingProduct.GetAllBookingProductsByFilter(ctx, sqlxRepo.GetAllBookingProductsByFilterParams{
		BookingID: bookingID,
		Limit:     &req.Limit,
		Offset:    &req.Offset,
		Sort:      &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking products", err)
	}

	// Convert to response format
	itemsDTO := make([]adminBookingProductModel.GetAllItemData, len(items))
	for i, item := range items {
		itemsDTO[i] = adminBookingProductModel.GetAllItemData{
			ID:   utils.FormatID(item.ProductID),
			Name: item.ProductName,
			Brand: adminBookingProductModel.GetAllItemProductBrandData{
				ID:   utils.FormatID(item.BrandID),
				Name: item.BrandName,
			},
			Category: adminBookingProductModel.GetAllItemProductCategoryData{
				ID:   utils.FormatID(item.CategoryID),
				Name: item.CategoryName,
			},
			CreatedAt: utils.PgTimestamptzToTimeString(item.CreatedAt),
		}
	}

	return &adminBookingProductModel.GetAllResponse{
		Total: total,
		Items: itemsDTO,
	}, nil
}
