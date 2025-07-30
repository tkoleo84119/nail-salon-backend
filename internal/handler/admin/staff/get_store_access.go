package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
)

type GetStaffStoreAccessHandler struct {
	service adminStaffService.GetStaffStoreAccessServiceInterface
}

func NewGetStaffStoreAccessHandler(service adminStaffService.GetStaffStoreAccessServiceInterface) *GetStaffStoreAccessHandler {
	return &GetStaffStoreAccessHandler{
		service: service,
	}
}

func (h *GetStaffStoreAccessHandler) GetStaffStoreAccess(c *gin.Context) {
	// Get staff ID from path parameter
	staffID := c.Param("staffId")
	if staffID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"staffId": "staffId 為必填項目"})
		return
	}

	// Service layer call
	response, err := h.service.GetStaffStoreAccess(c.Request.Context(), staffID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response - directly return the response data structure
	c.JSON(http.StatusOK, response)
}