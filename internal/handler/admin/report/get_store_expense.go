package adminReport

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminReportModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminReportService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/report"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreExpense struct {
	service adminReportService.GetStoreExpenseInterface
}

func NewGetStoreExpense(service adminReportService.GetStoreExpenseInterface) *GetStoreExpense {
	return &GetStoreExpense{
		service: service,
	}
}

func (h *GetStoreExpense) GetStoreExpense(c *gin.Context) {
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

	var req adminReportModel.GetStoreExpenseRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	startDate, err := utils.DateStringToTime(req.StartDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"startDate": "startDate 日期格式錯誤",
		})
		return
	}
	endDate, err := utils.DateStringToTime(req.EndDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"endDate": "endDate 日期格式錯誤",
		})
		return
	}

	staff, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIDs := make([]int64, len(staff.StoreList))
	for i, store := range staff.StoreList {
		storeIDs[i] = store.ID
	}

	parsedReq := adminReportModel.GetStoreExpenseParsedRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	response, err := h.service.GetStoreExpense(c.Request.Context(), storeID, parsedReq, staff.Role, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
