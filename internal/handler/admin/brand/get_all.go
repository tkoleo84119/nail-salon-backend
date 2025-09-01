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

type GetAll struct {
	service adminBrandService.GetAllInterface
}

func NewGetAll(service adminBrandService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Parse and validate query parameters
	var req adminBrandModel.GetAllRequest
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

	parsedReq := adminBrandModel.GetAllParsedRequest{
		Name:     req.Name,
		IsActive: req.IsActive,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
