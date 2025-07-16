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

type UpdateStylistHandler struct {
	updateStylistService stylistService.UpdateStylistServiceInterface
}

func NewUpdateStylistHandler(updateStylistService stylistService.UpdateStylistServiceInterface) *UpdateStylistHandler {
	return &UpdateStylistHandler{
		updateStylistService: updateStylistService,
	}
}

func (h *UpdateStylistHandler) UpdateStylist(c *gin.Context) {
	var req stylist.UpdateStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if utils.IsValidationError(err) {
			validationErrors := utils.ExtractValidationErrors(err)
			errorCodes.RespondWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		} else {
			fieldErrors := map[string]string{"request": "JSON格式錯誤"}
			errorCodes.RespondWithError(c, errorCodes.ValJsonFormat, fieldErrors)
		}
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

	response, err := h.updateStylistService.UpdateStylist(c.Request.Context(), req, staffUserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}