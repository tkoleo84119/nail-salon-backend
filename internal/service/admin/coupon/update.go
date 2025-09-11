package adminCoupon

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewUpdate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) UpdateInterface {
	return &Update{
		queries: queries,
		repo:    repo,
	}
}

func (s *Update) Update(ctx context.Context, couponID int64, req adminCouponModel.UpdateRequest) (*adminCouponModel.UpdateResponse, error) {
	// Ensure coupon exists
	exists, err := s.queries.CheckCouponExists(ctx, couponID)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check coupon existence", err)
	}
	if !exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNotFound)
	}

	// Name uniqueness excluding self
	if req.Name != nil {
		exists, err := s.queries.CheckCouponNameExistsExcluding(ctx, dbgen.CheckCouponNameExistsExcludingParams{
			ID:   couponID,
			Name: *req.Name,
		})
		if err != nil {
			return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check coupon name uniqueness", err)
		}
		if exists {
			return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNameAlreadyExists)
		}
	}

	// Perform partial update
	if err := s.repo.Coupon.UpdateCoupon(ctx, couponID, sqlxRepo.UpdateCouponParams{
		Name:     req.Name,
		IsActive: req.IsActive,
		Note:     req.Note,
	}); err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to update coupon", err)
	}

	response := &adminCouponModel.UpdateResponse{
		ID: utils.FormatID(couponID),
	}
	return response, nil
}
