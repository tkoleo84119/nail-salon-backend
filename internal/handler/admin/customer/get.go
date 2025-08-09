package adminCustomer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCustomerService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	service adminCustomerService.GetInterface
}

func NewGet(service adminCustomerService.GetInterface) *Get {
	return &Get{
		service: service,
	}
}

func (h *Get) Get(c *gin.Context) {
	// Get customer ID from path
	customerID := c.Param("customerId")
	if customerID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"customerId": "customerId 為必填項目",
		})
		return
	}
	parsedCustomerID, err := utils.ParseID(customerID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"customerId": "customerId 類型轉換失敗",
		})
		return
	}

	// Call service
	response, err := h.service.Get(c.Request.Context(), parsedCustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
