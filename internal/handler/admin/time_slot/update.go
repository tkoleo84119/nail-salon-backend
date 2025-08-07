package adminTimeSlot

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminTimeSlotModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminTimeSlotService.UpdateInterface
}

func NewUpdate(service adminTimeSlotService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	// Get path parameters
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"scheduleId": "scheduleId 為必填項目",
		})
		return
	}
	parsedScheduleID, err := utils.ParseID(scheduleID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"scheduleId": "scheduleId 類型轉換失敗",
		})
		return
	}

	timeSlotID := c.Param("timeSlotId")
	if timeSlotID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"timeSlotId": "timeSlotId 為必填項目",
		})
		return
	}
	parsedTimeSlotID, err := utils.ParseID(timeSlotID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"timeSlotId": "timeSlotId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminTimeSlotModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdate() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	updaterStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		updaterStoreIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.Update(c.Request.Context(), parsedScheduleID, parsedTimeSlotID, req, staffContext.UserID, staffContext.Role, updaterStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
