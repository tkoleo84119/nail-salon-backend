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

type DeleteSchedulesBulkHandler struct {
	service adminScheduleService.DeleteSchedulesBulkServiceInterface
}

func NewDeleteSchedulesBulkHandler(service adminScheduleService.DeleteSchedulesBulkServiceInterface) *DeleteSchedulesBulkHandler {
	return &DeleteSchedulesBulkHandler{
		service: service,
	}
}

func (h *DeleteSchedulesBulkHandler) DeleteSchedulesBulk(c *gin.Context) {
	// Parse and validate request
	var req adminScheduleModel.DeleteSchedulesBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
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
