package adminAuth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminAuthModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdatePassword struct {
	service adminAuthService.UpdatePasswordInterface
}

// NewUpdatePassword creates a new update password handler
func NewUpdatePassword(service adminAuthService.UpdatePasswordInterface) *UpdatePassword {
	return &UpdatePassword{
		service: service,
	}
}

// UpdatePassword handles the update password endpoint
func (h *UpdatePassword) UpdatePassword(c *gin.Context) {
	var req adminAuthModel.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Trim spaces
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	if req.OldPassword != nil {
		trimmed := strings.TrimSpace(*req.OldPassword)
		req.OldPassword = &trimmed
	}

	parseStaffId, err := utils.ParseID(req.StaffId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"staffId": "staffId 類型轉換失敗",
		})
		return
	}

	parsedReq := adminAuthModel.UpdatePasswordParsedRequest{
		StaffId:     parseStaffId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service layer
	response, err := h.service.UpdatePassword(c.Request.Context(), parsedReq, staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
