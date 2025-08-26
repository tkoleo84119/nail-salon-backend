package adminCustomerCoupon

import (
	"context"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries dbgen.Querier
}

func NewCreate(queries dbgen.Querier) *Create {
	return &Create{
		queries: queries,
	}
}

func (s *Create) Create(ctx context.Context, req adminCustomerCouponModel.CreateParsedRequest) (*adminCustomerCouponModel.CreateResponse, error) {
	// check customer and coupon existence
	exists, err := s.queries.CheckCustomerExistsByID(ctx, req.CustomerId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check customer existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerNotFound)
	}
	couponExists, err := s.queries.CheckCouponExists(ctx, req.CouponId)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check coupon existence", err)
	}
	if !couponExists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNotFound)
	}

	validFrom, validTo := s.setValidFromAndTo(req.Period)
	id := utils.GenerateID()

	err = s.queries.CreateCustomerCoupon(ctx, dbgen.CreateCustomerCouponParams{
		ID:         id,
		CustomerID: req.CustomerId,
		CouponID:   req.CouponId,
		ValidFrom:  utils.TimeToPgTimestamptz(validFrom),
		ValidTo:    utils.TimeToPgTimestamptz(validTo),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer coupon", err)
	}

	// response
	return &adminCustomerCouponModel.CreateResponse{
		ID: utils.FormatID(id),
	}, nil
}

// setValidFromAndTo 根據 period 設定 valid_from 與 valid_to
func (s *Create) setValidFromAndTo(period string) (time.Time, time.Time) {
	validFrom := time.Now()
	var validTo time.Time

	switch period {
	case "1month":
		validTo = validFrom.AddDate(0, 1, 0)
	case "3months":
		validTo = validFrom.AddDate(0, 3, 0)
	case "6months":
		validTo = validFrom.AddDate(0, 6, 0)
	case "1year":
		validTo = validFrom.AddDate(1, 0, 0)
	case "unlimited":
		validTo = time.Time{}
	default:
		validTo = time.Time{}
	}
	return validFrom, validTo
}
