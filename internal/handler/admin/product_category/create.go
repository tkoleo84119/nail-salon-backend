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

type Create struct {
	service adminProductCategoryService.CreateInterface
}

func NewCreate(service adminProductCategoryService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	var req adminProductCategoryModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name
	req.Name = strings.TrimSpace(req.Name)

	response, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
