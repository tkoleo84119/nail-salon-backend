package adminBookingProduct

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingProductModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking_product"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type BulkDelete struct {
	queries *dbgen.Queries
}

func NewBulkDelete(queries *dbgen.Queries) BulkDeleteInterface {
	return &BulkDelete{
		queries: queries,
	}
}

func (s *BulkDelete) BulkDelete(ctx context.Context, storeID int64, bookingID int64, req adminBookingProductModel.BulkDeleteParsedRequest, role string, staffStoreIDs []int64) (*adminBookingProductModel.BulkDeleteResponse, error) {
	if err := utils.CheckStoreAccess(storeID, staffStoreIDs, role); err != nil {
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

	// get all existing booking products
	existingProductIds, err := s.queries.GetAllBookingProductIdsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking products", err)
	}

	// filter products that actually exist in booking_products
	needDeleteProductIds := make([]int64, 0)
	for _, productID := range req.ProductIds {
		if slices.Contains(existingProductIds, productID) {
			needDeleteProductIds = append(needDeleteProductIds, productID)
		}
	}

	// if no products to delete, return empty result
	if len(needDeleteProductIds) == 0 {
		return &adminBookingProductModel.BulkDeleteResponse{
			Deleted: []string{},
		}, nil
	}

	// delete booking products
	err = s.queries.BulkDeleteBookingProducts(ctx, dbgen.BulkDeleteBookingProductsParams{
		BookingID: bookingID,
		Column2:   needDeleteProductIds,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete booking products", err)
	}

	deletedProductIds := make([]string, len(needDeleteProductIds))
	for i, productID := range needDeleteProductIds {
		deletedProductIds[i] = utils.FormatID(productID)
	}

	return &adminBookingProductModel.BulkDeleteResponse{
		Deleted: deletedProductIds,
	}, nil
}
