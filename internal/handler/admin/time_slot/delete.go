package adminTimeSlot

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Delete struct {
	service adminTimeSlotService.DeleteInterface
}

func NewDelete(service adminTimeSlotService.DeleteInterface) *Delete {
	return &Delete{
		service: service,
	}
}

func (h *Delete) Delete(c *gin.Context) {
	// Get path parameters
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"scheduleId": "scheduleId為必填項目",
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
			"timeSlotId": "timeSlotId為必填項目",
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

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	updaterID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"staffId": "staffId 類型轉換失敗",
		})
		return
	}

	updaterStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeID, err := utils.ParseID(store.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"storeId": "storeId 類型轉換失敗",
			})
			return
		}
		updaterStoreIDs[i] = storeID
	}

	// Call service
	response, err := h.service.Delete(c.Request.Context(), parsedScheduleID, parsedTimeSlotID, updaterID, staffContext.Role, updaterStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
