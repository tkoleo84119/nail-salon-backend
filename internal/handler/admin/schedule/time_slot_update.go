package adminSchedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotHandler struct {
	service adminScheduleService.UpdateTimeSlotServiceInterface
}

func NewUpdateTimeSlotHandler(service adminScheduleService.UpdateTimeSlotServiceInterface) *UpdateTimeSlotHandler {
	return &UpdateTimeSlotHandler{
		service: service,
	}
}

func (h *UpdateTimeSlotHandler) UpdateTimeSlot(c *gin.Context) {
	// Parse and validate request
	var req adminScheduleModel.UpdateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get path parameters
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"scheduleId": "scheduleId為必填項目",
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

	if !req.HasUpdate() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.UpdateTimeSlot(c.Request.Context(), scheduleID, timeSlotID, req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
