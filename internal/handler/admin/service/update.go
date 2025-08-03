package adminService

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateServiceHandler struct {
	service adminServiceService.UpdateServiceInterface
}

func NewUpdateServiceHandler(service adminServiceService.UpdateServiceInterface) *UpdateServiceHandler {
	return &UpdateServiceHandler{
		service: service,
	}
}

func (h *UpdateServiceHandler) UpdateService(c *gin.Context) {
	var req adminServiceModel.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get serviceId from path parameter
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"serviceId": "serviceId為必填項目",
		})
		return
	}
	parsedServiceID, err := utils.ParseID(serviceID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"serviceId": "serviceId 類型轉換失敗",
		})
		return
	}

	// Check if request has updates
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.UpdateService(c.Request.Context(), parsedServiceID, req, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
