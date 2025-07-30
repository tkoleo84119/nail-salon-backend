package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CustomerLineRegisterHandler struct {
	service authService.CustomerLineRegisterServiceInterface
}

func NewCustomerLineRegisterHandler(service authService.CustomerLineRegisterServiceInterface) *CustomerLineRegisterHandler {
	return &CustomerLineRegisterHandler{
		service: service,
	}
}

func (h *CustomerLineRegisterHandler) CustomerLineRegister(c *gin.Context) {
	var req authModel.CustomerLineRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Create login context
	loginCtx := authModel.CustomerLoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service
	response, err := h.service.CustomerLineRegister(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return successful response
	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}
