package adminTimeSlotTemplate

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminTimeSlotTemplateService.CreateInterface
}

func NewCreate(service adminTimeSlotTemplateService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	// Parse and validate request
	var req adminTimeSlotTemplateModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name, note
	req.Name = strings.TrimSpace(req.Name)
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.Create(c.Request.Context(), req, staffContext.UserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
