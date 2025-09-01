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

type GetAll struct {
	service adminProductCategoryService.GetAllInterface
}

func NewGetAll(service adminProductCategoryService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	var req adminProductCategoryModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	// set default value
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := adminProductCategoryModel.GetAllParsedRequest{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	}

	response, err := h.service.GetAll(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
