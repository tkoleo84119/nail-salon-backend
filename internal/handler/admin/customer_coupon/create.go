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

type Create struct {
	service adminCustomerCouponService.CreateInterface
}

func NewCreate(service adminCustomerCouponService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	var req adminCustomerCouponModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	customerID, err := utils.ParseID(req.CustomerId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"customerId": "customerId 類型轉換失敗",
		})
		return
	}

	couponID, err := utils.ParseID(req.CouponId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"couponId": "couponId 類型轉換失敗",
		})
		return
	}

	parsedReq := adminCustomerCouponModel.CreateParsedRequest{
		CustomerId: customerID,
		CouponId:   couponID,
		Period:     req.Period,
	}

	resp, err := h.service.Create(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(resp))
}
