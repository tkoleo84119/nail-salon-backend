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

func NewGet(queries *dbgen.Queries) *Get {
	return &Get{
		queries: queries,
	}
}

func (s *Get) Get(ctx context.Context, storeID, bookingID int64, role string, storeIds []int64) (*adminBookingModel.GetResponse, error) {
	// Check store access for staff (except SUPER_ADMIN)
	if role != common.RoleSuperAdmin {
		if err := utils.CheckStoreAccess(storeID, storeIds); err != nil {
			return nil, err
		}
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
		IsChatEnabled: utils.PgBoolToBool(booking.IsChatEnabled),
		Note:          utils.PgTextToString(booking.Note),
		Status:        booking.Status,
		CreatedAt:     utils.PgTimestamptzToTimeString(booking.CreatedAt),
		UpdatedAt:     utils.PgTimestamptzToTimeString(booking.UpdatedAt),
		Checkout:      nil, // default is nil
	}

	bookingDetails, err := s.queries.GetBookingDetailsByBookingID(ctx, bookingID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get booking details", err)
	}

	response.BookingDetails = make([]adminBookingModel.GetBookingDetailItem, len(bookingDetails))
	for i, detail := range bookingDetails {
		rawPrice := utils.PgNumericToFloat64(detail.Price)
		price := rawPrice

		if detail.DiscountRate.Valid {
			price = rawPrice * utils.PgNumericToFloat64(detail.DiscountRate)
		} else if detail.DiscountAmount.Valid {
			price = rawPrice - utils.PgNumericToFloat64(detail.DiscountAmount)
		}

		response.BookingDetails[i] = adminBookingModel.GetBookingDetailItem{
			ID: utils.FormatID(detail.ID),
			Service: adminBookingModel.GetService{
				ID:      utils.FormatID(detail.ServiceID),
				Name:    detail.ServiceName,
				IsAddon: utils.PgBoolToBool(detail.IsAddon),
			},
			RawPrice: int(rawPrice),
			Price:    int(price),
		}
	}

	if booking.Status == common.BookingStatusCompleted {
		checkout, err := s.queries.GetCheckoutByBookingID(ctx, bookingID)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get checkout", err)
		}

		response.Checkout = &adminBookingModel.GetCheckout{
			ID:            utils.FormatID(checkout.ID),
			PaymentMethod: checkout.PaymentMethod,
			TotalAmount:   int(utils.PgNumericToFloat64(checkout.TotalAmount)),
			FinalAmount:   int(utils.PgNumericToFloat64(checkout.FinalAmount)),
			PaidAmount:    int(utils.PgNumericToFloat64(checkout.PaidAmount)),
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
