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

type UpdateMyStaffHandler struct {
	service staffService.UpdateMyStaffServiceInterface
}

func NewUpdateMyStaffHandler(service staffService.UpdateMyStaffServiceInterface) *UpdateMyStaffHandler {
	return &UpdateMyStaffHandler{
		service: service,
	}
}

func (h *UpdateMyStaffHandler) UpdateMyStaff(c *gin.Context) {
	var req staff.UpdateMyStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if req.Email == nil {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
			"request": "至少需要提供一個欄位進行更新",
		})
		return
	}

	// Get staff user ID from JWT context
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	staffUserID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.UpdateMyStaff(c.Request.Context(), req, staffUserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
