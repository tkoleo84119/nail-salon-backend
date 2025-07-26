package schedule

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	scheduleService "github.com/tkoleo84119/nail-salon-backend/internal/service/schedule"
)

type GetTimeSlotHandler struct {
	service scheduleService.GetTimeSlotServiceInterface
}

func NewGetTimeSlotHandler(service scheduleService.GetTimeSlotServiceInterface) *GetTimeSlotHandler {
	return &GetTimeSlotHandler{
		service: service,
	}
}

func (h *GetTimeSlotHandler) GetTimeSlots(c *gin.Context) {
	// Path parameter validation
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"scheduleId": "scheduleId 為必填項目",
		})
		return
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetTimeSlotsBySchedule(c.Request.Context(), scheduleID, *customerContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response - return the array directly as per spec
	c.JSON(http.StatusOK, common.SuccessResponse(*response))
}
