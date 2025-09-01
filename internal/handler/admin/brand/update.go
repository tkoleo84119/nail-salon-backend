package adminBrand

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminBrandModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/brand"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminBrandService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/brand"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminBrandService.UpdateInterface
}

func NewUpdate(service adminBrandService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	brandIDStr := c.Param("brandId")
	if brandIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"brandId": "brandId 是必填項目",
		})
		return
	}
	brandID, err := utils.ParseID(brandIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"brandId": "brandId 類型轉換失敗",
		})
		return
	}

	var req adminBrandModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	response, err := h.service.Update(c.Request.Context(), brandID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
