package adminStockUsages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStockUsagesService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateFinish struct {
	service adminStockUsagesService.UpdateFinishInterface
}

func NewUpdateFinish(service adminStockUsagesService.UpdateFinishInterface) *UpdateFinish {
	return &UpdateFinish{
		service: service,
	}
}

func (h *UpdateFinish) UpdateFinish(c *gin.Context) {
	// Get path parameter - storeId
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 轉換類型失敗",
		})
		return
	}

	// Get path parameter - stockUsageId
	stockUsageID := c.Param("stockUsageId")
	if stockUsageID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"stockUsageId": "stockUsageId 為必填項目",
		})
		return
	}
	parsedStockUsageID, err := utils.ParseID(stockUsageID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stockUsageId": "stockUsageId 轉換類型失敗",
		})
		return
	}

	// Parse and validate JSON body
	var req adminStockUsagesModel.UpdateFinishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	usageEndedAt, err := utils.DateStringToTime(req.UsageEndedAt)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"usageEndedAt": "usageEndedAt 日期格式錯誤，應為 YYYY-MM-DD",
		})
		return
	}
	parsedReq := adminStockUsagesModel.UpdateFinishParsedRequest{
		UsageEndedAt: usageEndedAt,
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIds := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIds[i] = store.ID
	}

	// Service layer call
	response, err := h.service.UpdateFinish(c.Request.Context(), parsedStoreID, parsedStockUsageID, parsedReq, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
