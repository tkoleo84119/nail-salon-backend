package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteSchedulesBulkHandler struct {
	service scheduleService.DeleteSchedulesBulkServiceInterface
}

func NewDeleteSchedulesBulkHandler(service scheduleService.DeleteSchedulesBulkServiceInterface) *DeleteSchedulesBulkHandler {
	return &DeleteSchedulesBulkHandler{
		service: service,
	}
}

func (h *DeleteSchedulesBulkHandler) DeleteSchedulesBulk(c *gin.Context) {
	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Parse and validate request
	var req schedule.DeleteSchedulesBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Call service
	response, err := h.service.DeleteSchedulesBulk(c.Request.Context(), req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response with 204 No Content but with data
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}