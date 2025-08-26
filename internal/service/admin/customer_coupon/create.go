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

	// validate time
	now := time.Now()
	if req.ValidFrom.Before(now) {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponValidFromBeforeNow)
	}
	if req.ValidTo != nil {
		if req.ValidFrom.After(*req.ValidTo) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponValidFromAfterValidTo)
		}
		if req.ValidTo.Before(now) {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CustomerCouponValidToBeforeNow)
		}
	}

	id := utils.GenerateID()

	err = s.queries.CreateCustomerCoupon(ctx, dbgen.CreateCustomerCouponParams{
		ID:         id,
		CustomerID: req.CustomerId,
		CouponID:   req.CouponId,
		ValidFrom:  utils.TimeToPgTimestamptz(req.ValidFrom),
		ValidTo:    utils.TimePtrToPgTimestamptz(req.ValidTo),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create customer coupon", err)
	}

	// response
	return &adminCustomerCouponModel.CreateResponse{
		ID: utils.FormatID(id),
	}, nil
}
