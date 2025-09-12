package adminAuth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshToken struct {
	service adminAuthService.RefreshTokenInterface
	cfg     *config.Config
}

func NewRefreshToken(service adminAuthService.RefreshTokenInterface, cfg *config.Config) *RefreshToken {
	return &RefreshToken{
		service: service,
		cfg:     cfg,
	}
}

func (h *RefreshToken) RefreshToken(c *gin.Context) {
	// Read refresh token from HttpOnly cookie
	cookieToken, err := c.Cookie(h.cfg.Cookie.AdminRefreshName)
	if err != nil || strings.TrimSpace(cookieToken) == "" {
		errorCodes.AbortWithError(c, errorCodes.AuthRefreshTokenInvalid, map[string]string{
			"refreshToken": "refresh token 無效或已過期",
		})
		return
	}

	req := adminAuthModel.RefreshTokenRequest{RefreshToken: strings.TrimSpace(cookieToken)}

	refreshTokenCtx := adminAuthModel.RefreshTokenContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Service layer call
	response, err := h.service.RefreshToken(c.Request.Context(), req, refreshTokenCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Set the rotated refresh token back to cookie if present
	if strings.TrimSpace(response.RefreshToken) != "" {
		utils.SetAdminRefreshCookie(c, h.cfg.Cookie, response.RefreshToken)
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
