package adminCustomerCoupon

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCustomerCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer_coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	service adminCustomerCouponService.DeleteInterface
}

func NewDelete(service adminCustomerCouponService.DeleteInterface) *Delete {
	return &Delete{
		service: service,
	}
}

func (h *Delete) Delete(c *gin.Context) {
	// get customerCouponId from path parameter
	customerCouponIDStr := c.Param("customerCouponId")
	if customerCouponIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"customerCouponId": "customerCouponId 為必填項目",
		})
		return
	}
	customerCouponID, err := utils.ParseID(customerCouponIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"customerCouponId": "customerCouponId 類型轉換失敗",
		})
		return
	}

	// call service
	response, err := h.service.Delete(c.Request.Context(), customerCouponID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
