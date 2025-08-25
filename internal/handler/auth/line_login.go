package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineLogin struct {
	service authService.LineLoginInterface
}

// NewCustomerLineLoginHandler creates a new LINE login handler
func NewLineLogin(service authService.LineLoginInterface) *LineLogin {
	return &LineLogin{
		service: service,
	}
}

// CustomerLineLogin handles the customer LINE login endpoint
func (h *LineLogin) LineLogin(c *gin.Context) {
	// Input JSON Validation
	var req authModel.LineLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim idToken
	req.IdToken = strings.TrimSpace(req.IdToken)

	// Extract login context
	loginCtx := authModel.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service layer
	response, err := h.service.LineLogin(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
