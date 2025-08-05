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

type Update struct {
	service adminTimeSlotTemplateItemService.UpdateServiceInterface
}

func NewUpdate(service adminTimeSlotTemplateItemService.UpdateServiceInterface) *Update {
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

	itemID := c.Param("itemId")
	if itemID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"itemId": "itemId 為必填項目",
		})
		return
	}
	parsedItemID, err := utils.ParseID(itemID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"itemId": "itemId 類型轉換失敗",
		})
		return
	}

	var req adminTimeSlotTemplateItemModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	response, err := h.service.Update(c.Request.Context(), parsedTemplateID, parsedItemID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
