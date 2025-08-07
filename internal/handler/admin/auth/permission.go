package adminAuth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAuthService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/auth"
)

type Permission struct {
	service adminAuthService.PermissionInterface
}

func NewPermission(service adminAuthService.PermissionInterface) *Permission {
	return &Permission{
		service: service,
	}
}

func (h *Permission) GetPermission(c *gin.Context) {
	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Get permission data
	permission, err := h.service.Permission(c.Request.Context(), staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(permission))
}
