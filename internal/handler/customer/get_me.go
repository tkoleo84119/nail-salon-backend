package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
)

type GetMe struct {
	service *customerService.GetMe
}

func NewGetMe(service *customerService.GetMe) *GetMe {
	return &GetMe{
		service: service,
	}
}

func (h *GetMe) GetMe(c *gin.Context) {
	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetMe(c.Request.Context(), customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
