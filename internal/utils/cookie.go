package utils

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/tkoleo84119/nail-salon-backend/internal/config"
)

func sameSiteMode(s string) http.SameSite {
    switch strings.ToLower(strings.TrimSpace(s)) {
    case "lax":
        return http.SameSiteLaxMode
    case "strict":
        return http.SameSiteStrictMode
    case "none":
        return http.SameSiteNoneMode
    default:
        return http.SameSiteDefaultMode
    }
}

// SetAdminRefreshCookie sets the admin refresh token cookie based on config.Cookie settings.
func SetAdminRefreshCookie(c *gin.Context, cfg config.CookieConfig, token string) {
    maxAge := cfg.AdminRefreshMaxAgeDays * 24 * 60 * 60
    // configure SameSite before SetCookie
    c.SetSameSite(sameSiteMode(cfg.SameSite))
    c.SetCookie(cfg.AdminRefreshName, token, maxAge, cfg.Path, cfg.Domain, cfg.Secure, true)
}

// ClearAdminRefreshCookie clears the admin refresh token cookie (MaxAge=-1)
func ClearAdminRefreshCookie(c *gin.Context, cfg config.CookieConfig) {
    c.SetSameSite(sameSiteMode(cfg.SameSite))
    c.SetCookie(cfg.AdminRefreshName, "", -1, cfg.Path, cfg.Domain, cfg.Secure, true)
}

// SetCustomerRefreshCookie sets the customer refresh token cookie based on config.Cookie settings.
func SetCustomerRefreshCookie(c *gin.Context, cfg config.CookieConfig, token string) {
    maxAge := cfg.CustomerRefreshMaxAgeDays * 24 * 60 * 60
    c.SetSameSite(sameSiteMode(cfg.SameSite))
    c.SetCookie(cfg.CustomerRefreshName, token, maxAge, cfg.Path, cfg.Domain, cfg.Secure, true)
}

// ClearCustomerRefreshCookie clears the customer refresh token cookie (MaxAge=-1)
func ClearCustomerRefreshCookie(c *gin.Context, cfg config.CookieConfig) {
    c.SetSameSite(sameSiteMode(cfg.SameSite))
    c.SetCookie(cfg.CustomerRefreshName, "", -1, cfg.Path, cfg.Domain, cfg.Secure, true)
}
