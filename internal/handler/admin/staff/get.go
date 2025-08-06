package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	service adminStaffService.GetInterface
}

func NewGet(service adminStaffService.GetInterface) *Get {
	return &Get{
		service: service,
	}
}

func (h *Get) Get(c *gin.Context) {
	// Get staff ID from path parameter
	staffID := c.Param("staffId")
	if staffID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{"staffId": "staffId 為必填項目"})
		return
	}
	parsedStaffID, err := utils.ParseID(staffID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"staffId": "staffId 類型轉換失敗"})
		return
	}

	// Service layer call
	response, err := h.service.Get(c.Request.Context(), parsedStaffID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
