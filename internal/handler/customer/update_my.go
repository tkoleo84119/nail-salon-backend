package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	customerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyCustomerHandler struct {
	service customerService.UpdateMyCustomerServiceInterface
}

func NewUpdateMyCustomerHandler(service customerService.UpdateMyCustomerServiceInterface) *UpdateMyCustomerHandler {
	return &UpdateMyCustomerHandler{
		service: service,
	}
}

// UpdateMyCustomer handles PATCH /api/customers/me
func (h *UpdateMyCustomerHandler) UpdateMyCustomer(c *gin.Context) {
	var req customerModel.UpdateMyCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	customerID := customerContext.CustomerID
	result, err := h.service.UpdateMyCustomer(c.Request.Context(), customerID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(result))
}
