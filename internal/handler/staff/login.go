package staff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type LoginHandler struct {
	loginService staffService.LoginServiceInterface
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(loginService staffService.LoginServiceInterface) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

// Login handles the staff login endpoint
func (h *LoginHandler) Login(c *gin.Context) {
	var req staff.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle validation errors
		if utils.IsValidationError(err) {
			validationErrors := utils.ExtractValidationErrors(err)
			errorCodes.RespondWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		} else {
			// Handle JSON parsing errors
			fieldErrors := map[string]string{"request": "JSON格式錯誤"}
			errorCodes.RespondWithError(c, errorCodes.ValJsonFormat, fieldErrors)
		}
		return
	}

	// Extract login context
	loginCtx := staff.LoginContext{
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
	}

	// Call service layer
	response, err := h.loginService.Login(c.Request.Context(), req, loginCtx)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
