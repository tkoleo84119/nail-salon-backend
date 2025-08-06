package adminAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Logout struct {
	service adminAuthService.LogoutInterface
}

// NewStaffLogoutHandler creates a new logout handler
func NewLogout(service adminAuthService.LogoutInterface) *Logout {
	return &Logout{
		service: service,
	}
}

// StaffLogout handles the staff logout endpoint
func (h *Logout) Logout(c *gin.Context) {
	var req adminAuthModel.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Call the logout service
	response, err := h.service.Logout(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}