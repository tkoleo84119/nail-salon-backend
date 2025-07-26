package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleModel "github.com/tkoleo84119/nail-salon-backend/internal/model/schedule"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type ScheduleHandler struct {
	service scheduleService.ScheduleServiceInterface
}

func NewScheduleHandler(service scheduleService.ScheduleServiceInterface) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

func (h *ScheduleHandler) GetStoreSchedules(c *gin.Context) {
	// Path parameter validation
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"storeId": "門市ID為必填項目",
		})
		return
	}

	stylistID := c.Param("stylistId")
	if stylistID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"stylistId": "美甲師ID為必填項目",
		})
		return
	}

	// Query parameter validation
	var req scheduleModel.GetStoreSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetStoreSchedules(c.Request.Context(), storeID, stylistID, req, *customerContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
