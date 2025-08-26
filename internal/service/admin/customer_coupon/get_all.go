package adminCustomerCoupon

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
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

func (s *GetAll) GetAll(ctx context.Context, req adminCustomerCouponModel.GetAllParsedRequest) (*adminCustomerCouponModel.GetAllResponse, error) {
	total, results, err := s.repo.CustomerCoupon.GetAllCustomerCouponsByFilter(ctx, sqlxRepo.GetAllCustomerCouponsByFilterParams{
		CustomerID: req.CustomerId,
		CouponID:   req.CouponId,
		IsUsed:     req.IsUsed,
		Limit:      &req.Limit,
		Offset:     &req.Offset,
		Sort:       &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get customer coupons", err)
	}

	// get customer and coupon information
	customerIds := make([]int64, len(results))
	couponIds := make([]int64, len(results))
	for i, r := range results {
		customerIds[i] = r.CustomerID
		couponIds[i] = r.CouponID
	}
	customers, err := s.queries.GetCustomerByIDs(ctx, customerIds)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get customers", err)
	}
	coupons, err := s.queries.GetCouponByIDs(ctx, couponIds)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get coupons", err)
	}

	// create customer and coupon map
	customerMap := make(map[int64]dbgen.GetCustomerByIDsRow)
	for _, c := range customers {
		if _, ok := customerMap[c.ID]; !ok {
			customerMap[c.ID] = c
		}
	}
	couponMap := make(map[int64]dbgen.GetCouponByIDsRow)
	for _, c := range coupons {
		if _, ok := couponMap[c.ID]; !ok {
			couponMap[c.ID] = c
		}
	}

	items := make([]adminCustomerCouponModel.GetAllCustomerCouponItem, len(results))
	for i, r := range results {
		customer, ok := customerMap[r.CustomerID]
		if !ok {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Customer not found", nil)
		}
		coupon, ok := couponMap[r.CouponID]
		if !ok {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Coupon not found", nil)
		}

		items[i] = adminCustomerCouponModel.GetAllCustomerCouponItem{
			ID: utils.FormatID(r.ID),
			Customer: adminCustomerCouponModel.GetAllItemCustomerDTO{
				ID:       utils.FormatID(r.CustomerID),
				Name:     customer.Name,
				LineName: utils.PgTextToString(customer.LineName),
				Phone:    customer.Phone,
			},
			Coupon: adminCustomerCouponModel.GetAllItemCouponDTO{
				ID:             utils.FormatID(r.CouponID),
				DisplayName:    coupon.DisplayName,
				Code:           coupon.Code,
				DiscountRate:   utils.PgNumericToFloat64(coupon.DiscountRate),
				DiscountAmount: int64(utils.PgNumericToFloat64(coupon.DiscountAmount)),
				IsActive:       utils.PgBoolToBool(coupon.IsActive),
			},
			ValidFrom: utils.PgTimestamptzToTimeString(r.ValidFrom),
			ValidTo:   utils.PgTimestamptzToTimeString(r.ValidTo),
			IsUsed:    utils.PgBoolToBool(r.IsUsed),
			UsedAt:    utils.PgTimestamptzToTimeString(r.UsedAt),
			CreatedAt: utils.PgTimestamptzToTimeString(r.CreatedAt),
			UpdatedAt: utils.PgTimestamptzToTimeString(r.UpdatedAt),
		}
	}

	resp := &adminCustomerCouponModel.GetAllResponse{
		Total: total,
		Items: items,
	}
	return resp, nil
}
