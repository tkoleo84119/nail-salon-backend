package adminAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshToken struct {
	service adminAuthService.RefreshTokenInterface
}

func NewRefreshToken(service adminAuthService.RefreshTokenInterface) *RefreshToken {
	return &RefreshToken{
		service: service,
	}
}

func (h *RefreshToken) RefreshToken(c *gin.Context) {
	// Input JSON validation
	var req adminAuthModel.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
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
