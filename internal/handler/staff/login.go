package staff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
			c.JSON(http.StatusBadRequest, common.ValidationErrorResponse(validationErrors))
		} else {
			// Handle JSON parsing errors
			errors := map[string]string{"request": "JSON格式錯誤"}
			c.JSON(http.StatusBadRequest, common.ErrorResponse("請求錯誤", errors))
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
		// For security, don't expose internal errors
		if err.Error() == "invalid credentials" {
			errors := map[string]string{"credentials": "帳號或密碼錯誤"}
			c.JSON(http.StatusUnauthorized, common.ErrorResponse("認證失敗", errors))
		} else {
			errors := map[string]string{"server": "伺服器內部錯誤"}
			c.JSON(http.StatusInternalServerError, common.ErrorResponse("系統錯誤", errors))
		}
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
