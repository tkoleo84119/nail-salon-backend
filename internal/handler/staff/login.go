package staff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
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
		c.JSON(http.StatusBadRequest, staff.ErrorResponse{Error: "invalid request"})
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
			c.JSON(http.StatusUnauthorized, staff.ErrorResponse{Error: "invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, staff.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}
