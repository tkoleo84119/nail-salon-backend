package adminProductCategory

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminProductCategoryModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/product_category"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminProductCategoryService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/product_category"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminProductCategoryService.UpdateInterface
}

func NewUpdate(service adminProductCategoryService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
	productCategoryIDStr := c.Param("productCategoryId")
	if productCategoryIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"productCategoryId": "productCategoryId 是必填項目",
		})
		return
	}
	productCategoryID, err := utils.ParseID(productCategoryIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"productCategoryId": "productCategoryId 類型轉換失敗",
		})
		return
	}

	var req adminProductCategoryModel.UpdateRequest
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

	response, err := h.service.Update(c.Request.Context(), productCategoryID, req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
