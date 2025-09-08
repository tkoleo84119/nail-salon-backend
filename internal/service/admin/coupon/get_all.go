package adminCoupon

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	repo *sqlxRepo.Repositories
}

func NewGetAll(repo *sqlxRepo.Repositories) *GetAll {
	return &GetAll{repo: repo}
}

func (s *GetAll) GetAll(ctx context.Context, req adminCouponModel.GetAllParsedRequest) (*adminCouponModel.GetAllResponse, error) {
	total, results, err := s.repo.Coupon.GetAllCouponsByFilter(ctx, sqlxRepo.GetAllCouponsByFilterParams{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: req.IsActive,
		Limit:    &req.Limit,
		Offset:   &req.Offset,
		Sort:     &req.Sort,
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "Failed to get coupon list", err)
	}

	items := make([]adminCouponModel.GetAllCouponItemDTO, len(results))
	for i, result := range results {
		discountRate, err := utils.PgNumericToFloat64(result.DiscountRate)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount rate to float64", err)
		}
		discountAmount, err := utils.PgNumericToInt64(result.DiscountAmount)
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.ValTypeConversionFailed, "failed to convert discount amount to int64", err)
		}

		items[i] = adminCouponModel.GetAllCouponItemDTO{
			ID:             utils.FormatID(result.ID),
			Name:           result.Name,
			DisplayName:    result.DisplayName,
			Code:           result.Code,
			DiscountRate:   discountRate,
			DiscountAmount: discountAmount,
			IsActive:       utils.PgBoolToBool(result.IsActive),
			Note:           utils.PgTextToString(result.Note),
			CreatedAt:      utils.PgTimestamptzToTimeString(result.CreatedAt),
			UpdatedAt:      utils.PgTimestamptzToTimeString(result.UpdatedAt),
		}
	}

	response := &adminCouponModel.GetAllResponse{
		Total: total,
		Items: items,
	}
	return response, nil
}
