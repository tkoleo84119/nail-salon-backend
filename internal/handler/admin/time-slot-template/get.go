package adminTimeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
)

type GetTimeSlotTemplateHandler struct {
	service adminTimeSlotTemplateService.GetTimeSlotTemplateServiceInterface
}

func NewGetTimeSlotTemplateHandler(service adminTimeSlotTemplateService.GetTimeSlotTemplateServiceInterface) *GetTimeSlotTemplateHandler {
	return &GetTimeSlotTemplateHandler{
		service: service,
	}
}

func (h *GetTimeSlotTemplateHandler) GetTimeSlotTemplate(c *gin.Context) {
	// Get path parameter
	templateID := c.Param("templateId")
	if templateID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"templateId": "templateId 為必填項目"})
		return
	}

	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	template, err := h.service.GetTimeSlotTemplate(c.Request.Context(), templateID, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(template))
}
