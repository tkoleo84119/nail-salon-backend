package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminStaffService.UpdateInterface
}

func NewUpdate(service adminStaffService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	// Get target staff ID from path parameter
	staffId := c.Param("staffId")
	if staffId == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"staffId": "staffId 為必填項目",
		})
		return
	}
	parsedStaffId, err := utils.ParseID(staffId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"staffId": "staffId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminStaffModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Call service
	response, err := h.service.Update(c.Request.Context(), parsedStaffId, req, staffContext.UserID, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
