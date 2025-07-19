package timeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
)

type DeleteTimeSlotTemplateHandler struct {
	service timeSlotTemplateService.DeleteTimeSlotTemplateServiceInterface
}

func NewDeleteTimeSlotTemplateHandler(service timeSlotTemplateService.DeleteTimeSlotTemplateServiceInterface) *DeleteTimeSlotTemplateHandler {
	return &DeleteTimeSlotTemplateHandler{
		service: service,
	}
}

func (h *DeleteTimeSlotTemplateHandler) DeleteTimeSlotTemplate(c *gin.Context) {
	// Validate path parameter
	templateID := c.Param("templateId")
	if templateID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"templateId": "templateId為必填項目",
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
	response, err := h.service.DeleteTimeSlotTemplate(c.Request.Context(), templateID, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}