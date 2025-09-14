package adminCustomerCoupon

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	queries *dbgen.Queries
}

func NewDelete(queries *dbgen.Queries) DeleteInterface {
	return &Delete{
		queries: queries,
	}
}

func (s *Delete) Delete(ctx context.Context, customerCouponID int64) (*adminCustomerCouponModel.DeleteResponse, error) {
	// check customer coupon existence
	customerCoupon, err := s.queries.GetCustomerCouponForDelete(ctx, customerCouponID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponNotFound)
		}
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to get customer coupon", err)
	}

	// verify customer coupon is not used
	if customerCoupon.IsUsed.Valid && customerCoupon.IsUsed.Bool {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponAlreadyUsed)
	}

	// verify customer coupon is not expired
	if customerCoupon.ValidTo.Valid {
		validTo := customerCoupon.ValidTo.Time
		if time.Now().After(validTo) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponExpired)
		}
	}

	// delete customer coupon
	if err := s.queries.DeleteCustomerCoupon(ctx, customerCouponID); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to delete customer coupon", err)
	}

	// return delete result
	return &adminCustomerCouponModel.DeleteResponse{
		Deleted: utils.FormatID(customerCouponID),
	}, nil
}
