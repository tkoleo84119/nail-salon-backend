package adminService

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminServiceModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminServiceService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminServiceService.GetAllInterface
}

func NewGetAll(service adminServiceService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Parse and validate query parameters
	var req adminServiceModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name
	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	// default limit and offset
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := adminServiceModel.GetAllParsedRequest{
		Name:      req.Name,
		IsAddon:   req.IsAddon,
		IsActive:  req.IsActive,
		IsVisible: req.IsVisible,
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
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
