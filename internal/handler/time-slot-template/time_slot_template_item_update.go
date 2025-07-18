package timeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateTimeSlotTemplateItemHandler struct {
	service timeSlotTemplateService.UpdateTimeSlotTemplateItemServiceInterface
}

func NewUpdateTimeSlotTemplateItemHandler(service timeSlotTemplateService.UpdateTimeSlotTemplateItemServiceInterface) *UpdateTimeSlotTemplateItemHandler {
	return &UpdateTimeSlotTemplateItemHandler{
		service: service,
	}
}

func (h *UpdateTimeSlotTemplateItemHandler) UpdateTimeSlotTemplateItem(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		validationErrors := map[string]string{
			"templateId": "templateId為必填項目",
		}
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		validationErrors := map[string]string{
			"itemId": "itemId為必填項目",
		}
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	var req timeSlotTemplate.UpdateTimeSlotTemplateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	staffContext, exists := c.Get("staffContext")
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthTokenInvalid, nil)
		return
	}

	response, err := h.service.UpdateTimeSlotTemplateItem(c.Request.Context(), templateID, itemID, req, staffContext.(common.StaffContext))
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
