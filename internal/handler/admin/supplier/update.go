package adminSupplier

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminSupplierModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/supplier"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminSupplierService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/supplier"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminSupplierService.UpdateInterface
}

func NewUpdate(service adminSupplierService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	// Parse path parameter
	supplierIDStr := c.Param("supplierId")
	if supplierIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"supplierId": "supplierId 是必填項目",
		})
		return
	}
	supplierID, err := utils.ParseID(supplierIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"supplierId": "supplierId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request body
	var req adminSupplierModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{})
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	// Service layer call
	response, err := h.service.Update(c.Request.Context(), supplierID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
