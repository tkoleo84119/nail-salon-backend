package adminTimeSlotTemplate

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminTimeSlotTemplateService.UpdateInterface
}

func NewUpdate(service adminTimeSlotTemplateService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"templateId": "templateId 為必填項目",
		})
		return
	}
	parsedTemplateID, err := utils.ParseID(templateID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"templateId": "templateId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminTimeSlotTemplateModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get template ID from URL parameter
	if !req.HasUpdate() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// trim name, note
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	// Call service
	response, err := h.service.Update(c.Request.Context(), parsedTemplateID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
