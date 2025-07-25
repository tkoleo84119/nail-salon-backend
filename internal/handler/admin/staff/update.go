package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStaffHandler struct {
	service adminStaffService.UpdateStaffServiceInterface
}

func NewUpdateStaffHandler(service adminStaffService.UpdateStaffServiceInterface) *UpdateStaffHandler {
	return &UpdateStaffHandler{service: service}
}

func (h *UpdateStaffHandler) UpdateStaff(c *gin.Context) {
	// Parse and validate request
	var req adminStaffModel.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Get target staff ID from path parameter
	targetID := c.Param("staffId")
	if targetID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"staffId": "staffId為必填項目",
		})
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if !req.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Convert UserID to int64
	updaterID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.UpdateStaff(c.Request.Context(), targetID, req, updaterID, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
