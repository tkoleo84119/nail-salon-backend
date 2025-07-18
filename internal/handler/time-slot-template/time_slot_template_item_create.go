package timeSlotTemplate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotTemplate "github.com/tkoleo84119/nail-salon-backend/internal/model/time-slot-template"
	timeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateTimeSlotTemplateItemHandler struct {
	service timeSlotTemplateService.CreateTimeSlotTemplateItemServiceInterface
}

func NewCreateTimeSlotTemplateItemHandler(service timeSlotTemplateService.CreateTimeSlotTemplateItemServiceInterface) *CreateTimeSlotTemplateItemHandler {
	return &CreateTimeSlotTemplateItemHandler{
		service: service,
	}
}

func (h *CreateTimeSlotTemplateItemHandler) CreateTimeSlotTemplateItem(c *gin.Context) {
	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Get template ID from URL parameter
	templateID := c.Param("templateId")
	if templateID == "" {
		validationErrors := map[string]string{
			"templateId": "templateId為必填項目",
		}
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Parse and validate request
	var req timeSlotTemplate.CreateTimeSlotTemplateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Call service
	response, err := h.service.CreateTimeSlotTemplateItem(c.Request.Context(), templateID, req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
