package adminCheckout

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCheckoutModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/checkout"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries dbgen.Querier
	db      *pgxpool.Pool
	repo    *sqlxRepo.Repositories
}

type CouponInfo struct {
	ID             int64
	DiscountAmount *float64
	DiscountRate   *float64
}

func NewCreate(queries dbgen.Querier, repo *sqlxRepo.Repositories, db *pgxpool.Pool) *Create {
	return &Create{
		queries: queries,
		repo:    repo,
		db:      db,
	}
}

func (s *Create) Create(ctx context.Context, storeID int64, bookingID int64, req adminCheckoutModel.CreateParsedRequest, creatorID int64, storeIDs []int64) (*adminCheckoutModel.CreateResponse, error) {
	// check store access
	if err := utils.CheckStoreAccess(storeID, storeIDs); err != nil {
		return nil, err
	}

	// check booking exists
	booking, err := s.queries.GetBookingInfoByID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking status", err)
	}
	if booking.Status != "SCHEDULED" {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotCheckout)
	}
	if booking.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
	}

	// when customerCouponID is not nil, get coupon info
	couponInfo := CouponInfo{}
	if req.CustomerCouponID != nil {
		coupon, err := s.queries.GetCustomerCouponPriceInfoByID(ctx, *req.CustomerCouponID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer coupon price info", err)
		}

		if coupon.CustomerID != booking.CustomerID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponNotBelongToCustomer)
		}
		if coupon.IsUsed.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponAlreadyUsed)
		}
		if !coupon.IsActive.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNotActive)
		}

		if coupon.DiscountRate.Valid {
			discountRate := utils.PgNumericToFloat64(coupon.DiscountRate)
			couponInfo.DiscountRate = &discountRate
		}
		if coupon.DiscountAmount.Valid {
			discountAmount := utils.PgNumericToFloat64(coupon.DiscountAmount)
			couponInfo.DiscountAmount = &discountAmount
		}

		couponInfo.ID = coupon.CouponID
	}

	bookingDetailIDs := make([]int64, len(req.BookingDetails))
	for i, bookingDetail := range req.BookingDetails {
		bookingDetailIDs[i] = bookingDetail.ID
	}

	// check booking details exists
	count, err := s.queries.CountBookingDetailsByIDsAndBookingID(ctx, dbgen.CountBookingDetailsByIDsAndBookingIDParams{
		Column1:   bookingDetailIDs,
		BookingID: bookingID,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to count booking details", err)
	}
	if count != int64(len(bookingDetailIDs)) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingDetailNotFound)
	}

	updateBookingDetailPriceInfo, totalAmount, finalAmount, err := s.prepareBookingDetailPriceInfoAndValidate(req.BookingDetails, &couponInfo)
	if err != nil {
		return nil, err
	}

	totalAmountPg, err := utils.Float64ToPgNumeric(totalAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert total amount to pgtype.Numeric", err)
	}
	finalAmountPg, err := utils.Float64ToPgNumeric(finalAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert final amount to pgtype.Numeric", err)
	}
	paidAmountPg, err := utils.Int64ToPgNumeric(req.PaidAmount)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert paid amount to pgtype.Numeric", err)
	}

	newCheckout := dbgen.CreateCheckoutParams{
		ID:            utils.GenerateID(),
		BookingID:     bookingID,
		TotalAmount:   totalAmountPg,
		FinalAmount:   finalAmountPg,
		PaidAmount:    paidAmountPg,
		PaymentMethod: req.PaymentMethod,
		CheckoutUser:  utils.Int64ToPgInt8(creatorID),
	}

	if couponInfo.ID != 0 {
		newCheckout.CouponID = utils.Int64ToPgInt8(couponInfo.ID)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	err = qtx.CreateCheckout(ctx, newCheckout)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create checkout", err)
	}

	for _, updateBookingDetailPriceInfo := range updateBookingDetailPriceInfo {
		err = qtx.UpdateBookingDetailPriceInfo(ctx, updateBookingDetailPriceInfo)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking detail price info", err)
		}
	}

	err = qtx.UpdateBookingStatus(ctx, dbgen.UpdateBookingStatusParams{
		ID:     bookingID,
		Status: common.BookingStatusCompleted,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking status", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	return &adminCheckoutModel.CreateResponse{
		ID: utils.FormatID(bookingID),
	}, nil
}

func (s *Create) prepareBookingDetailPriceInfoAndValidate(
	passedBookingDetails []adminCheckoutModel.CreateBookingDetailParsed,
	couponInfo *CouponInfo,
) ([]dbgen.UpdateBookingDetailPriceInfoParams, float64, float64, error) {
	totalAmount := 0.0
	finalAmount := 0.0

	// if no booking details, return error
	if len(passedBookingDetails) == 0 {
		return nil, totalAmount, finalAmount,
			errorCodes.NewServiceErrorWithCode(errorCodes.BookingDetailNotFound)
	}

	bookingDetailPriceInfo := []dbgen.UpdateBookingDetailPriceInfoParams{}
	for _, bookingDetail := range passedBookingDetails {
		originalPrice := float64(bookingDetail.Price)
		discountedPrice := originalPrice

		discountRatePg := pgtype.Numeric{}
		discountAmountPg := pgtype.Numeric{}

		if couponInfo != nil {
			if couponInfo.DiscountRate != nil {
				discountedPrice = originalPrice * *couponInfo.DiscountRate

				ratePg, err := utils.Float64ToPgNumeric(*couponInfo.DiscountRate)
				if err != nil {
					return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount rate to pgtype.Numeric", err)
				}
				discountRatePg = ratePg
			} else if couponInfo.DiscountAmount != nil {
				discountedPrice = originalPrice - *couponInfo.DiscountAmount
				if discountedPrice < 0 {
					discountedPrice = 0 // not allow negative price
				}

				amountPg, err := utils.Float64ToPgNumeric(*couponInfo.DiscountAmount)
				if err != nil {
					return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount amount to pgtype.Numeric", err)
				}
				discountAmountPg = amountPg
			}
		}

		totalAmount += originalPrice
		finalAmount += discountedPrice

		discountedPricePg, err := utils.Float64ToPgNumeric(discountedPrice)
		if err != nil {
			return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to pgtype.Numeric", err)
		}

		bookingDetailPriceInfo = append(bookingDetailPriceInfo, dbgen.UpdateBookingDetailPriceInfoParams{
			ID:             bookingDetail.ID,
			Price:          discountedPricePg,
			DiscountRate:   discountRatePg,
			DiscountAmount: discountAmountPg,
		})
	}

	return bookingDetailPriceInfo, totalAmount, finalAmount, nil
}
