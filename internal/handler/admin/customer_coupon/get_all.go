package adminCustomerCoupon

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCustomerCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminCustomerCouponService.GetAllInterface
}

func NewGetAll(service adminCustomerCouponService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	var req adminCustomerCouponModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// pagination and sort
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	// parse ids if provided
	var customerIDPtr *int64
	var couponIDPtr *int64
	if req.CustomerId != nil && *req.CustomerId != "" {
		id, err := utils.ParseID(*req.CustomerId)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"customerId": "customerId 類型轉換失敗",
			})
			return
		}
		customerIDPtr = &id
	}
	if req.CouponId != nil && *req.CouponId != "" {
		id, err := utils.ParseID(*req.CouponId)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"couponId": "couponId 類型轉換失敗",
			})
			return
		}
		couponIDPtr = &id
	}

	parsedReq := adminCustomerCouponModel.GetAllParsedRequest{
		CustomerId: customerIDPtr,
		CouponId:   couponIDPtr,
		IsUsed:     req.IsUsed,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
	}

	resp, err := h.service.GetAll(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(resp))
}
