package adminCoupon

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCouponModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCouponService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/coupon"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminCouponService.UpdateInterface
}

func NewUpdate(service adminCouponService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	couponIDStr := c.Param("couponId")
	if couponIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"couponId": "couponId 為必填項目",
		})
		return
	}

	couponID, err := utils.ParseID(couponIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"couponId": "couponId 類型轉換失敗",
		})
		return
	}

	var req adminCouponModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// trim name, note
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	response, err := h.service.Update(c.Request.Context(), couponID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
