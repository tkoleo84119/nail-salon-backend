package adminStockUsages

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStockUsagesModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStockUsagesService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stock_usages"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminStockUsagesService.CreateInterface
}

func NewCreate(service adminStockUsagesService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	storeIDStr := c.Param("storeId")
	if storeIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	var req adminStockUsagesModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	productID, err := utils.ParseID(req.ProductID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"productId": "productId 類型轉換失敗",
		})
		return
	}

	usageStarted, err := utils.DateStringToTime(req.UsageStarted)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"usageStarted": "usageStarted 日期格式錯誤，應為 YYYY-MM-DD",
		})
		return
	}

	var expirationDate *time.Time
	if req.Expiration != nil {
		expiration, err := utils.DateStringToTime(*req.Expiration)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"expiration": "expiration 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		expirationDate = &expiration
	}

	quantity := 1
	if req.Quantity != nil {
		quantity = int(*req.Quantity)
	}

	parsedReq := adminStockUsagesModel.CreateParsedRequest{
		ProductID:    productID,
		Quantity:     quantity,
		Expiration:   expirationDate,
		UsageStarted: usageStarted,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	response, err := h.service.Create(c.Request.Context(), storeID, parsedReq, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
