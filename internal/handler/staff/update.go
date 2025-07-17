package staff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStaffHandler struct {
	service staffService.UpdateStaffServiceInterface
}

func NewUpdateStaffHandler(service staffService.UpdateStaffServiceInterface) *UpdateStaffHandler {
	return &UpdateStaffHandler{service: service}
}

func (h *UpdateStaffHandler) UpdateStaff(c *gin.Context) {
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Get target staff ID from path parameter
	targetID := c.Param("id")
	if targetID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"id": "id為必填項目",
		})
		return
	}

	// Parse and validate request
	var req staff.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if req.Role == nil && req.IsActive == nil {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
			"request": "至少需要提供一個欄位進行更新",
		})
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
