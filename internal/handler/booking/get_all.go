package booking

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	bookingModel "github.com/tkoleo84119/nail-salon-backend/internal/model/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	bookingService "github.com/tkoleo84119/nail-salon-backend/internal/service/booking"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service bookingService.GetAllInterface
}

func NewGetAll(service bookingService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

var validStatuses = map[string]struct{}{
	"SCHEDULED": {},
	"CANCELLED": {},
	"COMPLETED": {},
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Query parameter validation
	var req bookingModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set default value
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	// Parse status parameter if provided (supports comma-separated values)
	statuses := make([]string, 0)
	if req.Status != nil && *req.Status != "" {
		rowStatuses := strings.Split(*req.Status, ",")
		for _, status := range rowStatuses {
			if _, ok := validStatuses[strings.ToUpper(strings.TrimSpace(status))]; !ok {
				errorCodes.AbortWithError(c, errorCodes.ValFieldOneof, map[string]string{"status": "status 必須是 SCHEDULED, CANCELLED, COMPLETED 其中一個值"})
				return
			}

			statuses = append(statuses, strings.TrimSpace(status))
		}
	}

	parsedReq := bookingModel.GetAllParsedRequest{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Status: &statuses,
	}

	// Authentication context validation
	customerContext, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedReq, customerContext.CustomerID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
