package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type AcceptTerm struct {
	service authService.AcceptTermInterface
}

// NewAcceptTerm creates a new accept term handler
func NewAcceptTerm(service authService.AcceptTermInterface) *AcceptTerm {
	return &AcceptTerm{
		service: service,
	}
}

// AcceptTerm handles the customer accept terms endpoint
func (h *AcceptTerm) AcceptTerm(c *gin.Context) {
	var req authModel.AcceptTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	customer, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthCustomerFailed, nil)
		return
	}

	response, err := h.service.AcceptTerm(c.Request.Context(), req, customer.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
