package adminAuth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "github.com/tkoleo84119/nail-salon-backend/internal/config"
    errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
    adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
    "github.com/tkoleo84119/nail-salon-backend/internal/model/common"
    adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
    "github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Logout struct {
    service adminAuthService.LogoutInterface
    cfg     *config.Config
}

// NewStaffLogoutHandler creates a new logout handler
func NewLogout(service adminAuthService.LogoutInterface, cfg *config.Config) *Logout {
    return &Logout{
        service: service,
        cfg:     cfg,
    }
}

// StaffLogout handles the staff logout endpoint
func (h *Logout) Logout(c *gin.Context) {
    // Read refresh token from HttpOnly cookie
    cookieToken, _ := c.Cookie(h.cfg.Cookie.AdminRefreshName)
    req := adminAuthModel.LogoutRequest{RefreshToken: strings.TrimSpace(cookieToken)}

    // Call the logout service
    response, err := h.service.Logout(c.Request.Context(), req)
    if err != nil {
        errorCodes.RespondWithServiceError(c, err)
        return
    }

    // Clear the refresh token cookie
    utils.ClearAdminRefreshCookie(c, h.cfg.Cookie)

    c.JSON(http.StatusOK, common.SuccessResponse(response))
}
