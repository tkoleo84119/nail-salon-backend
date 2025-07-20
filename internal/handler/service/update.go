package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateServiceHandler struct {
	updateServiceService serviceService.UpdateServiceInterface
}

func NewUpdateServiceHandler(updateServiceService serviceService.UpdateServiceInterface) *UpdateServiceHandler {
	return &UpdateServiceHandler{
		updateServiceService: updateServiceService,
	}
}

func (h *UpdateServiceHandler) UpdateService(c *gin.Context) {
	var req service.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Get serviceId from path parameter
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"serviceId": "serviceId為必填項目",
		})
		return
	}

	// Check if request has updates
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
			"request": "至少需要提供一個欄位進行更新",
		})
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.updateServiceService.UpdateService(c.Request.Context(), serviceID, req, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
