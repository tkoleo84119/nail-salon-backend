package adminStylist

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStylistModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/stylist"
	adminStylistService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/stylist"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateMyStylistHandler struct {
	service adminStylistService.UpdateMyStylistServiceInterface
}

func NewUpdateMyStylistHandler(service adminStylistService.UpdateMyStylistServiceInterface) *UpdateMyStylistHandler {
	return &UpdateMyStylistHandler{
		service: service,
	}
}

func (h *UpdateMyStylistHandler) UpdateMyStylist(c *gin.Context) {
	var req adminStylistModel.UpdateMyStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Additional validation: ensure at least one field is provided for update
	if !req.HasUpdate() {
		errorCodes.RespondWithEmptyFieldError(c)
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

	response, err := h.service.UpdateMyStylist(c.Request.Context(), req, staffUserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
