package adminSchedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminScheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/schedule"
	adminScheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateSchedulesBulkHandler struct {
	service adminScheduleService.CreateSchedulesBulkServiceInterface
}

func NewCreateSchedulesBulkHandler(service adminScheduleService.CreateSchedulesBulkServiceInterface) *CreateSchedulesBulkHandler {
	return &CreateSchedulesBulkHandler{
		service: service,
	}
}

func (h *CreateSchedulesBulkHandler) CreateSchedulesBulk(c *gin.Context) {
	// Parse and validate request
	var req adminScheduleModel.CreateSchedulesBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.CreateSchedulesBulk(c.Request.Context(), req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
