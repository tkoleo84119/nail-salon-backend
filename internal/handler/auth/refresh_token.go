package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshTokenHandler struct {
	service authService.RefreshTokenServiceInterface
}

func NewRefreshTokenHandler(service authService.RefreshTokenServiceInterface) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		service: service,
	}
}

func (h *RefreshTokenHandler) RefreshToken(c *gin.Context) {
	// Input JSON validation
	var req authModel.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Service layer call
	response, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}