package adminBookingProduct

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking_product"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type BulkCreate struct {
	queries *dbgen.Queries
}

func NewBulkCreate(queries *dbgen.Queries) BulkCreateInterface {
	return &BulkCreate{
		queries: queries,
	}
}

func (s *BulkCreate) BulkCreate(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.BulkCreateParsedRequest, staffStoreIDs []int64) (*adminBookingProductModel.BulkCreateResponse, error) {
	if err := utils.CheckStoreAccess(storeID, staffStoreIDs); err != nil {
		return nil, err
	}

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
	// Verify booking is in COMPLETED status
	if booking.Status != common.BookingStatusCompleted {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotAllowedToUpdate)
	}

	// Verify products exist and belong to the store
	productCount, err := s.queries.CountProductsByIDs(ctx, dbgen.CountProductsByIDsParams{
		Column1: req.ProductIds,
		StoreID: storeID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to count products", err)
	}
	if productCount != int64(len(req.ProductIds)) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ProductNotFound)
	}

	// get all exist booking products
	existingProductIds, err := s.queries.GetAllBookingProductIdsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking products", err)
	}

	// compare existing product ids and request product ids
	needCreateProductIds := make([]int64, 0)
	for _, productID := range req.ProductIds {
		if !slices.Contains(existingProductIds, productID) {
			needCreateProductIds = append(needCreateProductIds, productID)
		}
	}

	// create booking products
	now := time.Now()
	nowPg := utils.TimePtrToPgTimestamptz(&now)
	bookingProductsData := make([]dbgen.BulkCreateBookingProductsParams, len(needCreateProductIds))
	for i, productID := range needCreateProductIds {
		bookingProductsData[i] = dbgen.BulkCreateBookingProductsParams{
			BookingID: bookingID,
			ProductID: productID,
			CreatedAt: nowPg,
		}
	}

	_, err = s.queries.BulkCreateBookingProducts(ctx, bookingProductsData)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create booking products", err)
	}

	createdProductIds := make([]string, len(needCreateProductIds))
	for i, productID := range needCreateProductIds {
		createdProductIds[i] = utils.FormatID(productID)
	}

	return &adminBookingProductModel.BulkCreateResponse{
		Created: createdProductIds,
	}, nil
}
