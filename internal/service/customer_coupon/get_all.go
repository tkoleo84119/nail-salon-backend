package customerCoupon

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	customerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewGetAll(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{
		queries: queries,
		repo:    repo,
	}
}

func (s *GetAll) GetAll(ctx context.Context, customerID int64, req customerCouponModel.GetAllParsedRequest) (*customerCouponModel.GetAllResponse, error) {
	// set default sort
	if len(req.Sort) == 0 {
		req.Sort = []string{"isUsed", "validTo"}
	}

	total, results, err := s.repo.CustomerCoupon.GetAllCustomerCouponsByFilter(ctx, sqlxRepo.GetAllCustomerCouponsByFilterParams{
		CustomerID: &customerID,
		IsUsed:     req.IsUsed,
		Limit:      &req.Limit,
		Offset:     &req.Offset,
		Sort:       &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get customer coupons", err)
	}

	if total == 0 {
		return &customerCouponModel.GetAllResponse{
			Total: 0,
			Items: []customerCouponModel.GetAllCustomerCouponItem{},
		}, nil
	}

	couponIds := make([]int64, len(results))
	for i, r := range results {
		couponIds[i] = r.CouponID
	}
	coupons, err := s.queries.GetCouponByIDs(ctx, couponIds)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get coupons", err)
	}

	couponMap := make(map[int64]dbgen.GetCouponByIDsRow)
	for _, c := range coupons {
		if _, ok := couponMap[c.ID]; !ok {
			couponMap[c.ID] = c
		}
	}

	items := make([]customerCouponModel.GetAllCustomerCouponItem, len(results))
	for i, r := range results {
		coupon, ok := couponMap[r.CouponID]
		if !ok {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Coupon not found", nil)
		}

		items[i] = customerCouponModel.GetAllCustomerCouponItem{
			ID:        utils.FormatID(r.ID),
			ValidFrom: utils.PgTimestamptzToTimeString(r.ValidFrom),
			ValidTo:   utils.PgTimestamptzToTimeString(r.ValidTo),
			IsUsed:    utils.PgBoolToBool(r.IsUsed),
			UsedAt:    utils.PgTimestamptzToTimeString(r.UsedAt),
			CreatedAt: utils.PgTimestamptzToTimeString(r.CreatedAt),
			Coupon: customerCouponModel.GetAllItemCouponDTO{
				ID:             utils.FormatID(r.CouponID),
				DisplayName:    coupon.DisplayName,
				DiscountRate:   utils.PgNumericToFloat64(coupon.DiscountRate),
				DiscountAmount: int64(utils.PgNumericToFloat64(coupon.DiscountAmount)),
				IsActive:       utils.PgBoolToBool(coupon.IsActive),
			},
		}
	}

	return &customerCouponModel.GetAllResponse{
		Total: total,
		Items: items,
	}, nil
}
