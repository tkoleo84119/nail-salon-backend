package adminTimeSlotTemplateItem

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time_slot_template_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateItemService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time_slot_template_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminTimeSlotTemplateItemService.CreateInterface
}

func NewCreate(service adminTimeSlotTemplateItemService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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
	var req adminTimeSlotTemplateItemModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Call service
	response, err := h.service.Create(c.Request.Context(), parsedTemplateID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
