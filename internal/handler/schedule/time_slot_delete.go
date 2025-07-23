package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
)

type DeleteTimeSlotHandler struct {
	service scheduleService.DeleteTimeSlotServiceInterface
}

func NewDeleteTimeSlotHandler(service scheduleService.DeleteTimeSlotServiceInterface) *DeleteTimeSlotHandler {
	return &DeleteTimeSlotHandler{
		service: service,
	}
}

func (h *DeleteTimeSlotHandler) DeleteTimeSlot(c *gin.Context) {
	// Get path parameters
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"scheduleId": "scheduleId為必填項目",
		})
		return
	}

	timeSlotID := c.Param("timeSlotId")
	if timeSlotID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"timeSlotId": "timeSlotId為必填項目",
		})
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.DeleteTimeSlot(c.Request.Context(), scheduleID, timeSlotID, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
