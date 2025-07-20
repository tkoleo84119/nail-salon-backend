package customer

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/customer"
	customerService "github.com/tkoleo84119/nail-salon-backend/internal/service/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineLoginHandler struct {
	lineLoginService customerService.LineLoginServiceInterface
}

// NewLineLoginHandler creates a new LINE login handler
func NewLineLoginHandler(lineLoginService customerService.LineLoginServiceInterface) *LineLoginHandler {
	return &LineLoginHandler{
		lineLoginService: lineLoginService,
	}
}

// LineLogin handles the customer LINE login endpoint
func (h *LineLoginHandler) LineLogin(c *gin.Context) {
	// Input JSON Validation
	var req customer.LineLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Extract login context
	loginCtx := customer.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service layer
	response, err := h.lineLoginService.LineLogin(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}