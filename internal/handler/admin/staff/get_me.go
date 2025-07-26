package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetMyStaffHandler struct {
	service adminStaffService.GetMyStaffServiceInterface
}

func NewGetMyStaffHandler(service adminStaffService.GetMyStaffServiceInterface) *GetMyStaffHandler {
	return &GetMyStaffHandler{
		service: service,
	}
}

func (h *GetMyStaffHandler) GetMyStaff(c *gin.Context) {
	// Get staff context from JWT middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthTokenInvalid, nil)
		return
	}

	// Parse staff user ID from context
	staffUserID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{"staffUserId": "invalid staff user ID"})
		return
	}

	// Service layer call
	response, err := h.service.GetMyStaff(c.Request.Context(), staffUserID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}