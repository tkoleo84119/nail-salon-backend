package stylist

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/stylist"
	stylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyStylistHandler struct {
	updateMyStylistService stylistService.UpdateMyStylistServiceInterface
}

func NewUpdateMyStylistHandler(updateMyStylistService stylistService.UpdateMyStylistServiceInterface) *UpdateMyStylistHandler {
	return &UpdateMyStylistHandler{
		updateMyStylistService: updateMyStylistService,
	}
}

func (h *UpdateMyStylistHandler) UpdateMyStylist(c *gin.Context) {
	var req stylist.UpdateMyStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if !req.HasUpdate() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
			"request": "至少需要提供一個欄位進行更新",
		})
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	staffUserID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.updateMyStylistService.UpdateMyStylist(c.Request.Context(), req, staffUserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
