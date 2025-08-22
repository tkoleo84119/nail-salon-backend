package adminSchedule

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminScheduleService.UpdateInterface
}

func NewUpdate(service adminScheduleService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	// Validate path parameters
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

	scheduleIDStr := c.Param("scheduleId")
	if scheduleIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"scheduleId": "scheduleId 為必填項目",
		})
		return
	}
	scheduleID, err := utils.ParseID(scheduleIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"scheduleId": "scheduleId 類型轉換失敗",
		})
		return
	}

	// Parse request body
	var req adminScheduleModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValJsonFormat, nil)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	parsedStylistID, err := utils.ParseID(req.StylistID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"stylistId": "stylistId 類型轉換失敗",
		})
		return
	}

	var workDate *time.Time
	if req.WorkDate != nil {
		parsedWorkDate, err := utils.DateStringToTime(*req.WorkDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"workDate": "workDate 類型轉換失敗",
			})
			return
		}
		workDate = &parsedWorkDate
	}

	parsedRequest := adminScheduleModel.UpdateParsedRequest{
		StylistID: parsedStylistID,
		WorkDate:  workDate,
		Note:      req.Note,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Extract store IDs from staff context
	storeIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIDs[i] = store.ID
	}

	// Call service to update schedule
	result, err := h.service.Update(c.Request.Context(), storeID, scheduleID, parsedRequest, staffContext.UserID, staffContext.Role, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	response := common.SuccessResponse(result)
	c.JSON(http.StatusOK, response)
}
