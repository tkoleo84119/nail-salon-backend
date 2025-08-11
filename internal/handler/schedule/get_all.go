package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service scheduleService.GetAllInterface
}

func NewGetAll(service scheduleService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Path parameter validation
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
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	stylistID := c.Param("stylistId")
	if stylistID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"stylistId": "stylistId 為必填項目",
		})
		return
	}
	parsedStylistID, err := utils.ParseID(stylistID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stylistId": "stylistId 類型轉換失敗",
		})
		return
	}

	// Query parameter validation
	var req scheduleModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"startDate": "startDate 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)",
		})
		return
	}
	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"endDate": "endDate 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD)",
		})
		return
	}

	parsedReq := scheduleModel.GetAllParsedRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedStylistID, parsedReq, customerContext.IsBlacklisted)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
