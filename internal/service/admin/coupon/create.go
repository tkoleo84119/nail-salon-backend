package adminCoupon

import (
	"context"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlc/dbgen"
	sqlxRepo "github.com/tkoleo84119/nail-salon-backend/internal/repository/sqlx"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	queries *dbgen.Queries
	repo    *sqlxRepo.Repositories
}

func NewCreate(queries *dbgen.Queries, repo *sqlxRepo.Repositories) *Create {
	return &Create{
		queries: queries,
		repo:    repo,
	}
}

func (s *Create) Create(ctx context.Context, req adminCouponModel.CreateRequest) (*adminCouponModel.CreateResponse, error) {
	// Validate DiscountRate and DiscountAmount pass at least one
	if err := s.ValidateDiscountRateAndDiscountAmount(req); err != nil {
		return nil, err
	}

	// Check if coupon name already exists
	exists, err := s.queries.CheckCouponNameExists(ctx, req.Name)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check coupon name existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponNameAlreadyExists)
	}

	// Check if coupon code already exists
	exists, err = s.queries.CheckCouponCodeExists(ctx, req.Code)
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to check coupon code existence", err)
	}
	if exists {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.CouponCodeAlreadyExists)
	}

	// Convert to pgtype.Numeric
	discountRate, err := utils.Float64PtrToPgNumeric(req.DiscountRate)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValTypeConversionFailed)
	}
	discountAmount, err := utils.Int64PtrToPgNumeric(req.DiscountAmount)
	if err != nil {
		return nil, errorCodes.NewServiceErrorWithCode(errorCodes.ValTypeConversionFailed)
	}

	couponID := utils.GenerateID()
	isActive := true

	err = s.queries.CreateCoupon(ctx, dbgen.CreateCouponParams{
		ID:             couponID,
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Code:           req.Code,
		DiscountRate:   discountRate,
		DiscountAmount: discountAmount,
		IsActive:       utils.BoolPtrToPgBool(&isActive),
		Note:           utils.StringPtrToPgText(req.Note, true),
	})
	if err != nil {
		return nil, errorCodes.NewServiceError(errorCodes.SysDatabaseError, "failed to create coupon", err)
	}

	response := &adminCouponModel.CreateResponse{
		ID: utils.FormatID(couponID),
	}

	return response, nil
}

// ValidateDiscountRateAndDiscountAmount validates the discount rate and discount amount
func (s *Create) ValidateDiscountRateAndDiscountAmount(req adminCouponModel.CreateRequest) error {
	if req.DiscountRate == nil && req.DiscountAmount == nil {
		return errorCodes.NewServiceErrorWithCode(errorCodes.CouponDiscountRequired)
	}
	if req.DiscountRate != nil && req.DiscountAmount != nil {
		return errorCodes.NewServiceErrorWithCode(errorCodes.CouponDiscountExclusive)
	}

	return nil
}
