package customer

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

// LineRegisterServiceInterface defines the interface for LINE register service
type LineRegisterServiceInterface interface {
	LineRegister(ctx context.Context, req customer.LineRegisterRequest, loginCtx customer.LoginContext) (*customer.LineRegisterResponse, error)
}

type LineRegisterHandler struct {
	service LineRegisterServiceInterface
}

func NewLineRegisterHandler(service LineRegisterServiceInterface) *LineRegisterHandler {
	return &LineRegisterHandler{
		service: service,
	}
}

func (h *LineRegisterHandler) LineRegister(c *gin.Context) {
	var req customer.LineRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Create login context
	loginCtx := customer.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service
	response, err := h.service.LineRegister(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return successful response
	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}
