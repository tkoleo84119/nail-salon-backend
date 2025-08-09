package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LineRegister struct {
	service authService.LineRegisterInterface
}

func NewLineRegister(service authService.LineRegisterInterface) *LineRegister {
	return &LineRegister{
		service: service,
	}
}

func (h *LineRegister) LineRegister(c *gin.Context) {
	var req authModel.LineRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Create login context
	loginCtx := authModel.LoginContext{
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
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
