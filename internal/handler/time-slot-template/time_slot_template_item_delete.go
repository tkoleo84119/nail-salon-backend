package timeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
)

type DeleteTimeSlotTemplateItemHandler struct {
	service timeSlotTemplateService.DeleteTimeSlotTemplateItemServiceInterface
}

func NewDeleteTimeSlotTemplateItemHandler(service timeSlotTemplateService.DeleteTimeSlotTemplateItemServiceInterface) *DeleteTimeSlotTemplateItemHandler {
	return &DeleteTimeSlotTemplateItemHandler{
		service: service,
	}
}

func (h *DeleteTimeSlotTemplateItemHandler) DeleteTimeSlotTemplateItem(c *gin.Context) {
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

	staffContext, exists := c.Get("staffContext")
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthTokenInvalid, nil)
		return
	}

	response, err := h.service.DeleteTimeSlotTemplateItem(c.Request.Context(), templateID, itemID, staffContext.(common.StaffContext))
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
