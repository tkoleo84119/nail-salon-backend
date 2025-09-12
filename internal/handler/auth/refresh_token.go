package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tkoleo84119/nail-salon-backend/internal/config"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	authModel "github.com/tkoleo84119/nail-salon-backend/internal/model/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	authService "github.com/tkoleo84119/nail-salon-backend/internal/service/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type RefreshToken struct {
	service authService.RefreshTokenInterface
	cfg     *config.Config
}

func NewRefreshToken(service authService.RefreshTokenInterface, cfg *config.Config) *RefreshToken {
	return &RefreshToken{
		service: service,
		cfg:     cfg,
	}
}

func (h *RefreshToken) RefreshToken(c *gin.Context) {
	// Read refresh token from HttpOnly cookie
	cookieToken, err := c.Cookie(h.cfg.Cookie.CustomerRefreshName)
	if err != nil || strings.TrimSpace(cookieToken) == "" {
		errorCodes.AbortWithError(c, errorCodes.AuthRefreshTokenInvalid, map[string]string{
			"refreshToken": "refresh token 無效或已過期",
		})
		return
	}

	req := authModel.RefreshTokenRequest{RefreshToken: strings.TrimSpace(cookieToken)}

	refreshTokenCtx := authModel.RefreshTokenContext{
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

	// Set rotated refresh token back to cookie if present
	if strings.TrimSpace(response.RefreshToken) != "" {
		utils.SetCustomerRefreshCookie(c, h.cfg.Cookie, response.RefreshToken)
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
