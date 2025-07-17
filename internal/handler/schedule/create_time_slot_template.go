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

type CreateTimeSlotTemplateHandler struct {
	service scheduleService.CreateTimeSlotTemplateServiceInterface
}

func NewCreateTimeSlotTemplateHandler(service scheduleService.CreateTimeSlotTemplateServiceInterface) *CreateTimeSlotTemplateHandler {
	return &CreateTimeSlotTemplateHandler{
		service: service,
	}
}

func (h *CreateTimeSlotTemplateHandler) CreateTimeSlotTemplate(c *gin.Context) {
	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Parse and validate request
	var req scheduleModel.CreateTimeSlotTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Call service
	response, err := h.service.CreateTimeSlotTemplate(c.Request.Context(), req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}