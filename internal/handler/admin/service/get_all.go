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

type GetServiceListHandler struct {
	service adminServiceService.GetServiceListServiceInterface
}

func NewGetServiceListHandler(service adminServiceService.GetServiceListServiceInterface) *GetServiceListHandler {
	return &GetServiceListHandler{
		service: service,
	}
}

func (h *GetServiceListHandler) GetServiceList(c *gin.Context) {
	// Parse and validate query parameters
	var req adminServiceModel.GetServiceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// default limit and offset
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	sort := []string{}
	if req.Sort != nil {
		sort = strings.Split(*req.Sort, ",")
	}

	parsedReq := adminServiceModel.GetServiceListParsedRequest{
		Name:      req.Name,
		IsAddon:   req.IsAddon,
		IsActive:  req.IsActive,
		IsVisible: req.IsVisible,
		Limit:     limit,
		Offset:    offset,
		Sort:      sort,
	}

	// Service layer call
	response, err := h.service.GetServiceList(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
