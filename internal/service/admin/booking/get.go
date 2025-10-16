package adminBooking

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	queries *dbgen.Queries
}

func NewGet(queries *dbgen.Queries) GetInterface {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID, bookingID int64, role string, storeIds []int64) (*adminBookingModel.GetResponse, error) {
	// Check store access for staff
	if err := utils.CheckStoreAccess(storeID, storeIds, role); err != nil {
		return nil, err
	}

	booking, err := s.queries.GetBookingDetailByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get booking", err)
	}
	if booking.StoreID != storeID {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.BookingNotBelongToStore)
	}

	response := adminBookingModel.GetResponse{
		ID: utils.FormatID(booking.ID),
		Customer: adminBookingModel.GetCustomer{
			ID:   utils.FormatID(booking.CustomerID),
			Name: booking.CustomerName,
		},
		Stylist: adminBookingModel.GetStylist{
			ID:   utils.FormatID(booking.StylistID),
			Name: utils.PgTextToString(booking.StylistName),
		},
		TimeSlot: adminBookingModel.GetTimeSlot{
			ID:        utils.FormatID(booking.TimeSlotID),
			WorkDate:  utils.PgDateToDateString(booking.WorkDate),
			StartTime: utils.PgTimeToTimeString(booking.StartTime),
			EndTime:   utils.PgTimeToTimeString(booking.EndTime),
		},
		IsChatEnabled:  utils.PgBoolToBool(booking.IsChatEnabled),
		Note:           utils.PgTextToString(booking.Note),
		ActualDuration: utils.PgInt4ToInt32Ptr(booking.ActualDuration),
		Status:         booking.Status,
		CreatedAt:      utils.PgTimestamptzToTimeString(booking.CreatedAt),
		UpdatedAt:      utils.PgTimestamptzToTimeString(booking.UpdatedAt),
		Checkout:       nil, // default is nil
	}

	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get booking details", err)
	}

	response.BookingDetails = make([]adminBookingModel.GetBookingDetailItem, len(bookingDetails))
	for i, detail := range bookingDetails {
		rawPrice, err := utils.PgNumericToFloat64(detail.Price)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert price to float64", err)
		}
		price := rawPrice

		if detail.DiscountRate.Valid {
			discountRate, err := utils.PgNumericToFloat64(detail.DiscountRate)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount rate to float64", err)
			}

			price = rawPrice * discountRate
		} else if detail.DiscountAmount.Valid {
			discountAmount, err := utils.PgNumericToFloat64(detail.DiscountAmount)
			if err != nil {
				return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount amount to float64", err)
			}

			price = rawPrice - discountAmount
		}

		response.BookingDetails[i] = adminBookingModel.GetBookingDetailItem{
			ID: utils.FormatID(detail.ID),
			Service: adminBookingModel.GetService{
				ID:      utils.FormatID(detail.ServiceID),
				Name:    detail.ServiceName,
				IsAddon: utils.PgBoolToBool(detail.IsAddon),
			},
			RawPrice: rawPrice,
			Price:    price,
		}
	}

	if booking.Status == common.BookingStatusCompleted {
		checkout, err := s.queries.GetCheckoutByBookingID(ctx, bookingID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get checkout", err)
		}

		totalAmount, err := utils.PgNumericToInt64(checkout.TotalAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert total amount to int64", err)
		}
		finalAmount, err := utils.PgNumericToInt64(checkout.FinalAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert final amount to int64", err)
		}
		paidAmount, err := utils.PgNumericToInt64(checkout.PaidAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert paid amount to int64", err)
		}

		response.Checkout = &adminBookingModel.GetCheckout{
			ID:            utils.FormatID(checkout.ID),
			PaymentMethod: checkout.PaymentMethod,
			TotalAmount:   totalAmount,
			FinalAmount:   finalAmount,
			PaidAmount:    paidAmount,
			CheckoutUser:  utils.PgTextToString(checkout.CheckoutUser),
			Coupon:        nil, // default is nil
		}

		if checkout.CouponID.Valid {
			response.Checkout.Coupon = &adminBookingModel.GetCoupon{
				ID:   utils.PgInt8ToIDString(checkout.CouponID),
				Name: utils.PgTextToString(checkout.CouponName),
				Code: utils.PgTextToString(checkout.CouponCode),
			}
		}
	}

	return &response, nil
}
