package adminService

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetServiceHandler struct {
	service adminServiceService.GetServiceServiceInterface
}

func NewGetServiceHandler(service adminServiceService.GetServiceServiceInterface) *GetServiceHandler {
	return &GetServiceHandler{
		service: service,
	}
}

func (h *GetServiceHandler) GetService(c *gin.Context) {
	// Get service ID from path parameter
	serviceID := c.Param("serviceId")
	if serviceID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"serviceId": "serviceId 為必填項目"})
		return
	}
	parsedServiceID, err := utils.ParseID(serviceID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"serviceId": "serviceId 類型轉換失敗"})
		return
	}

	// Service layer call
	response, err := h.service.GetService(c.Request.Context(), parsedServiceID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
