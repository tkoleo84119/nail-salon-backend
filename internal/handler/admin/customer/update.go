package adminCustomer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCustomerService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminCustomerService.UpdateInterface
}

func NewUpdate(service adminCustomerService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
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

	var req adminCustomerModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	response, err := h.service.Update(c.Request.Context(), parsedCustomerID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
