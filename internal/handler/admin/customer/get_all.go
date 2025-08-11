package adminCustomer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminCustomerModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminCustomerService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/customer"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminCustomerService.GetAllInterface
}

func NewGetAll(service adminCustomerService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Bind query parameters
	var req adminCustomerModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set limit and offset
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := adminCustomerModel.GetAllParsedRequest{
		Name:          req.Name,
		LineName:      req.LineName,
		Phone:         req.Phone,
		Level:         req.Level,
		IsBlacklisted: req.IsBlacklisted,
		MinPastDays:   req.MinPastDays,
		Limit:         limit,
		Offset:        offset,
		Sort:          sort,
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
