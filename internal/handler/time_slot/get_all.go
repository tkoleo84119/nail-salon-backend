package timeSlot

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	timeSlotService "github.com/tkoleo84119/nail-salon-backend/internal/service/time_slot"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service timeSlotService.GetAllInterface
}

func NewGetAll(service timeSlotService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Path parameter validation
	scheduleID := c.Param("scheduleId")
	if scheduleID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"scheduleId": "scheduleId 為必填項目",
		})
		return
	}
	parsedScheduleID, err := utils.ParseID(scheduleID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"scheduleId": "scheduleId 類型轉換失敗",
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
	response, err := h.service.GetAll(c.Request.Context(), parsedScheduleID, customerContext.IsBlacklisted)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response - return the array directly as per spec
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
