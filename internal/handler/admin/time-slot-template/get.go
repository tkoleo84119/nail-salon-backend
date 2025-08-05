package adminTimeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
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
	parsedTemplateID, err := utils.ParseID(templateID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"templateId": "templateId 類型轉換失敗"})
		return
	}

	// Call service
	template, err := h.service.GetTimeSlotTemplate(c.Request.Context(), parsedTemplateID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(template))
}
