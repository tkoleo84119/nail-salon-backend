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

type Login struct {
	service adminAuthService.LoginInterface
	cfg     *config.Config
}

// NewStaffLoginHandler creates a new login handler
func NewLogin(service adminAuthService.LoginInterface, cfg *config.Config) *Login {
	return &Login{
		service: service,
		cfg:     cfg,
	}
}

// StaffLogin handles the staff login endpoint
func (h *Login) Login(c *gin.Context) {
	var req adminAuthModel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim username and password
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// Extract login context
	loginCtx := adminAuthModel.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service layer
	response, err := h.service.Login(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Set refresh token in HttpOnly cookie and do not return it in JSON
	if strings.TrimSpace(response.RefreshToken) != "" {
		utils.SetAdminRefreshCookie(c, h.cfg.Cookie, response.RefreshToken)
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
