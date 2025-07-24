package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
)

type GetMyCustomerHandler struct {
	service *customerService.GetMyCustomerService
}

func NewGetMyCustomerHandler(service *customerService.GetMyCustomerService) *GetMyCustomerHandler {
	return &GetMyCustomerHandler{
		service: service,
	}
}

func (h *GetMyCustomerHandler) GetMyCustomer(c *gin.Context) {
	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetMyCustomer(c.Request.Context(), *customerContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}