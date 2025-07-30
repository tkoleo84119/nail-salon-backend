package adminAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type StaffRefreshTokenHandler struct {
	service adminAuthService.StaffRefreshTokenServiceInterface
}

func NewStaffRefreshTokenHandler(service adminAuthService.StaffRefreshTokenServiceInterface) *StaffRefreshTokenHandler {
	return &StaffRefreshTokenHandler{
		service: service,
	}
}

func (h *StaffRefreshTokenHandler) StaffRefreshToken(c *gin.Context) {
	// Input JSON validation
	var req adminAuthModel.StaffRefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Service layer call
	response, err := h.service.StaffRefreshToken(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
