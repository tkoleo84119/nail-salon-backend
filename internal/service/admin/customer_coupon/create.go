package adminCustomerCoupon

import (
	"context"
	"fmt"
	"time"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
}

func NewCreate(queries *dbgen.Queries) CreateInterface {
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

	validFrom, validTo, err := s.setValidFromAndTo(req.Period)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysInternalError, "failed to set valid from and to", err)
	}
	id := utils.GenerateID()

	err = s.queries.CreateCustomerCoupon(ctx, dbgen.CreateCustomerCouponParams{
		ID:         id,
		CustomerID: req.CustomerId,
		CouponID:   req.CouponId,
		ValidFrom:  utils.TimePtrToPgTimestamptz(&validFrom),
		ValidTo:    utils.TimePtrToPgTimestamptz(&validTo),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer coupon", err)
	}

	// response
	return &adminCustomerCouponModel.CreateResponse{
		ID: utils.FormatID(id),
	}, nil
}

// setValidFromAndTo set valid_from and valid_to based on period
func (s *Create) setValidFromAndTo(period string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}

	validFrom := time.Now().In(loc)
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
		return validFrom, time.Time{}, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s", period)
	}

	// set time to 23:59:59
	validTo = time.Date(
		validTo.Year(),
		validTo.Month(),
		validTo.Day(),
		23, 59, 59, 0,
		loc,
	)

	return validFrom, validTo, nil
}
