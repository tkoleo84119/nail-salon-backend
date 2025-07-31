package adminAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffLogoutHandler struct {
	service adminAuthService.StaffLogoutServiceInterface
}

// NewStaffLogoutHandler creates a new logout handler
func NewStaffLogoutHandler(service adminAuthService.StaffLogoutServiceInterface) *StaffLogoutHandler {
	return &StaffLogoutHandler{
		service: service,
	}
}

// StaffLogout handles the staff logout endpoint
func (h *StaffLogoutHandler) StaffLogout(c *gin.Context) {
	var req adminAuthModel.StaffLogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Call the logout service
	response, err := h.service.StaffLogout(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}