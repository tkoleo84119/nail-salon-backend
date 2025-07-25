package adminSchedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateTimeSlotHandler struct {
	service adminScheduleService.CreateTimeSlotServiceInterface
}

func NewCreateTimeSlotHandler(service adminScheduleService.CreateTimeSlotServiceInterface) *CreateTimeSlotHandler {
	return &CreateTimeSlotHandler{
		service: service,
	}
}

func (h *CreateTimeSlotHandler) CreateTimeSlot(c *gin.Context) {
	// Parse and validate request
	var req adminScheduleModel.CreateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Get schedule ID from path parameter
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"scheduleId": "scheduleId為必填項目",
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
	response, err := h.service.CreateTimeSlot(c.Request.Context(), scheduleID, req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
