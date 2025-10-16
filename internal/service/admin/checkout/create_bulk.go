package adminCheckout

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCheckoutModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/checkout"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/service/cache"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateBulk struct {
	queries     *dbgen.Queries
	db          *pgxpool.Pool
	repo        *sqlxRepo.Repositories
	activityLog cache.ActivityLogCacheInterface
}

type CouponInfo struct {
	ID             int64
	DiscountAmount *float64
	DiscountRate   *float64
	ApplyCount     int64
}

func NewCreateBulk(queries *dbgen.Queries, repo *sqlxRepo.Repositories, db *pgxpool.Pool, activityLog cache.ActivityLogCacheInterface) CreateBulkInterface {
	return &CreateBulk{
		queries:     queries,
		repo:        repo,
		db:          db,
		activityLog: activityLog,
	}
}

func (s *CreateBulk) CreateBulk(ctx context.Context, storeID int64, req adminCheckoutModel.CreateBulkParsedRequest, staffContext *common.StaffContext) (*adminCheckoutModel.CreateBulkResponse, error) {
	storeIDs := make([]int64, len(staffContext.StoreList))
	storeName := "門市"
	for i, store := range staffContext.StoreList {
		storeIDs[i] = store.ID
		if store.ID == storeID {
			storeName = store.Name
		}
	}

	// check store access
	if err := utils.CheckStoreAccess(storeID, storeIDs, staffContext.Role); err != nil {
		return nil, err
	}

	bookingDetailMap := make(map[int64]dbgen.GetBookingDetailPriceInfoByBookingIDRow)
	customerIDs := make([]int64, len(req.Checkouts))
	applyCount := int64(0)
	for i, checkout := range req.Checkouts {
		// check booking exists
		booking, err := s.queries.GetBookingInfoWithDateByID(ctx, checkout.BookingID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking status", err)
		}
		if booking.Status != common.BookingStatusScheduled {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingStatusNotCheckout)
		}
		if booking.StoreID != storeID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
		}
		// check booking not in future
		if err := CheckBookingNotInFuture(booking.WorkDate, booking.StartTime); err != nil {
			return nil, err
		}

		bookingDetailPriceInfo, err := s.queries.GetBookingDetailPriceInfoByBookingID(ctx, checkout.BookingID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get booking detail price info", err)
		}
		if len(bookingDetailPriceInfo) != len(checkout.Details) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingDetailNotFound)
		}

		// map booking detail price info
		for _, bookingDetailPriceInfo := range bookingDetailPriceInfo {
			bookingDetailMap[bookingDetailPriceInfo.ID] = bookingDetailPriceInfo
		}

		customerIDs[i] = booking.CustomerID
		applyCount += checkout.ApplyCount
	}

	// check customerID is same in all checkouts
	customerID := customerIDs[0]
	if len(customerIDs) > 1 {
		for i := range customerIDs {
			if i != 0 && customerID != customerIDs[i] {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingWithMultipleCustomersNotAllowedToCheckout)
			}
		}
	}

	// when customerCouponID is not nil, get coupon info
	couponInfo := CouponInfo{
		ApplyCount: applyCount,
	}
	if req.CustomerCouponID != nil {
		coupon, err := s.queries.GetCustomerCouponPriceInfoByID(ctx, *req.CustomerCouponID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer coupon price info", err)
		}

		if coupon.CustomerID != customerID {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponNotBelongToCustomer)
		}
		if coupon.IsUsed.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponAlreadyUsed)
		}
		if !coupon.IsActive.Bool {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNotActive)
		}
		if coupon.ValidTo.Valid && coupon.ValidTo.Time.Before(time.Now()) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponExpired)
		}

		if coupon.DiscountRate.Valid {
			discountRate, err := utils.PgNumericToFloat64(coupon.DiscountRate)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount rate to float64", err)
			}

			couponInfo.DiscountRate = &discountRate
		}

		if coupon.DiscountAmount.Valid {
			discountAmount, err := utils.PgNumericToFloat64(coupon.DiscountAmount)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount amount to float64", err)
			}

			// check amount / apply count is not allow decimal
			if int64(discountAmount)%applyCount != 0 {
				return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponDiscountAmountNotDivisibleByApplyCount)
			}

			couponInfo.DiscountAmount = &discountAmount
		}

		couponInfo.ID = coupon.CouponID
	}

	newCheckouts, needUpdateBookingDetailPriceInfos, bookingIDs, err := s.prepareCheckoutAndUpdateBookingDetailData(req.PaymentMethod, req.Checkouts, bookingDetailMap, staffContext.UserID, &couponInfo)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := dbgen.New(tx)

	_, err = qtx.BulkCreateCheckout(ctx, newCheckouts)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create checkout", err)
	}

	for _, updateBookingDetailPriceInfo := range needUpdateBookingDetailPriceInfos {
		err = qtx.UpdateBookingDetailPriceInfo(ctx, updateBookingDetailPriceInfo)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update booking detail price info", err)
		}
	}

	err = qtx.UpdateBookingsStatus(ctx, dbgen.UpdateBookingsStatusParams{
		Column1: bookingIDs,
		Status:  common.BookingStatusCompleted,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update bookings status", err)
	}

	if req.CustomerCouponID != nil {
		err = qtx.UpdateCustomerCouponUsed(ctx, *req.CustomerCouponID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer coupon used", err)
		}
	}

	// update customer last visit at
	err = qtx.UpdateCustomerLastVisitAt(ctx, customerID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update customer last visit at", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to commit transaction", err)
	}

	// Log activity
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		customer, err := s.queries.GetCustomerByID(logCtx, customerID)
		if err == nil {
			if err := s.activityLog.LogAdminBookingCompleted(logCtx, staffContext.Username, customer.Name, utils.PgTextToString(customer.LineName), len(newCheckouts), storeName); err != nil {
				log.Printf("failed to log admin booking completed activity: %v", err)
			}
		}
	}()

	ids := make([]string, len(newCheckouts))
	for i, newCheckout := range newCheckouts {
		ids[i] = utils.FormatID(newCheckout.ID)
	}

	return &adminCheckoutModel.CreateBulkResponse{
		IDs: ids,
	}, nil
}

// CheckBookingNotInFuture checks if the booking is in the future
func CheckBookingNotInFuture(workDatePg pgtype.Date, startTime pgtype.Time) error {
	workDate := workDatePg.Time
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to load location", err)
	}

	startT, err := utils.PgTimeToTime(startTime)
	if err != nil {
		return errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert time", err)
	}

	bookingDateTime := time.Date(
		workDate.Year(),
		workDate.Month(),
		workDate.Day(),
		startT.Hour(),
		startT.Minute(),
		startT.Second(),
		startT.Nanosecond(),
		loc,
	)

	now := time.Now().In(loc)
	if bookingDateTime.After(now) {
		return errorCodes.NewServiceErrorWithCode(errorCodes.BookingInFutureNotAllowedToCheckout)
	}

	return nil
}

func (s *CreateBulk) prepareCheckoutAndUpdateBookingDetailData(
	paymentMethod string,
	passedBookings []adminCheckoutModel.CreateBulkParsedCheckoutItems,
	bookingDetailMap map[int64]dbgen.GetBookingDetailPriceInfoByBookingIDRow,
	creatorID int64,
	couponInfo *CouponInfo,
) ([]dbgen.BulkCreateCheckoutParams, []dbgen.UpdateBookingDetailPriceInfoParams, []int64, error) {
	newCheckouts := []dbgen.BulkCreateCheckoutParams{}
	needUpdateBookingDetailPriceInfos := []dbgen.UpdateBookingDetailPriceInfoParams{}
	bookingIDs := make([]int64, len(passedBookings))

	now := time.Now()
	nowPg := utils.TimePtrToPgTimestamptz(&now)

	for i, booking := range passedBookings {
		updateBookingDetailPriceInfos, totalAmount, finalAmount, err := s.prepareUpdateBookingDetail(booking.Details, bookingDetailMap, couponInfo)
		if err != nil {
			return nil, nil, nil, err
		}

		totalAmountPg, err := utils.Float64PtrToPgNumeric(&totalAmount)
		if err != nil {
			return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert total amount to pgtype.Numeric", err)
		}
		finalAmountPg, err := utils.Float64PtrToPgNumeric(&finalAmount)
		if err != nil {
			return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert final amount to pgtype.Numeric", err)
		}
		paidAmountPg, err := utils.Int64PtrToPgNumeric(&booking.PaidAmount)
		if err != nil {
			return nil, nil, nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert paid amount to pgtype.Numeric", err)
		}

		newCheckouts = append(newCheckouts, dbgen.BulkCreateCheckoutParams{
			ID:            utils.GenerateID(),
			BookingID:     booking.BookingID,
			TotalAmount:   totalAmountPg,
			FinalAmount:   finalAmountPg,
			PaidAmount:    paidAmountPg,
			PaymentMethod: paymentMethod,
			CouponID:      utils.Int64PtrToPgInt8(&couponInfo.ID),
			CheckoutUser:  utils.Int64PtrToPgInt8(&creatorID),
			CreatedAt:     nowPg,
			UpdatedAt:     nowPg,
		})

		bookingIDs[i] = booking.BookingID
		needUpdateBookingDetailPriceInfos = append(needUpdateBookingDetailPriceInfos, updateBookingDetailPriceInfos...)
	}

	return newCheckouts, needUpdateBookingDetailPriceInfos, bookingIDs, nil
}

func (s *CreateBulk) prepareUpdateBookingDetail(
	passedBookingDetails []adminCheckoutModel.CreateBulkParsedDetailItems,
	bookingDetailMap map[int64]dbgen.GetBookingDetailPriceInfoByBookingIDRow,
	couponInfo *CouponInfo,
) ([]dbgen.UpdateBookingDetailPriceInfoParams, float64, float64, error) {
	totalAmount := 0.0
	finalAmount := 0.0

	// if no booking details, return error
	if len(passedBookingDetails) == 0 {
		return nil, totalAmount, finalAmount, errorCodes.NewServiceErrorWithCode(errorCodes.BookingDetailNotFound)
	}

	bookingDetailPriceInfo := []dbgen.UpdateBookingDetailPriceInfoParams{}
	for _, bookingDetail := range passedBookingDetails {
		rawBookingDetail, ok := bookingDetailMap[bookingDetail.ID]
		if !ok {
			return nil, totalAmount, finalAmount, errorCodes.NewServiceErrorWithCode(errorCodes.BookingDetailNotFound)
		}

		originalPrice := float64(bookingDetail.Price)
		discountedPrice := originalPrice

		discountRatePg := pgtype.Numeric{Valid: false}
		discountAmountPg := pgtype.Numeric{Valid: false}

		// if coupon info exists and booking detail use coupon, apply coupon
		if couponInfo != nil && bookingDetail.UseCoupon {
			if couponInfo.DiscountRate != nil {
				discountedPrice = originalPrice * *couponInfo.DiscountRate

				ratePg, err := utils.Float64PtrToPgNumeric(couponInfo.DiscountRate)
				if err != nil {
					return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount rate to pgtype.Numeric", err)
				}
				discountRatePg = ratePg
			} else if couponInfo.DiscountAmount != nil {
				// when coupon is DiscountAmount, apply discount amount to each booking detail (need to divide by apply count)
				discountAmount := *couponInfo.DiscountAmount / float64(couponInfo.ApplyCount)
				discountedPrice = originalPrice - discountAmount

				if discountedPrice < 0 {
					discountedPrice = 0 // not allow negative price
				}

				amountPg, err := utils.Float64PtrToPgNumeric(&discountAmount)
				if err != nil {
					return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount amount to pgtype.Numeric", err)
				}
				discountAmountPg = amountPg
			}
		}

		totalAmount += originalPrice
		finalAmount += discountedPrice

		originalPricePg, err := utils.Float64PtrToPgNumeric(&originalPrice)
		if err != nil {
			return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to pgtype.Numeric", err)
		}

		rawOriginalPrice, err := utils.PgNumericToFloat64(rawBookingDetail.Price)
		if err != nil {
			return nil, totalAmount, finalAmount, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert original price to float64", err)
		}

		// if original price or discount rate or discount amount is changed, update booking detail price info
		if rawOriginalPrice != originalPrice || discountRatePg.Valid || discountAmountPg.Valid {
			bookingDetailPriceInfo = append(bookingDetailPriceInfo, dbgen.UpdateBookingDetailPriceInfoParams{
				ID:             bookingDetail.ID,
				Price:          originalPricePg, // store passed original price
				DiscountRate:   discountRatePg,
				DiscountAmount: discountAmountPg,
			})
		}
	}

	finalAmount = math.Round(finalAmount)

	return bookingDetailPriceInfo, totalAmount, finalAmount, nil
}
