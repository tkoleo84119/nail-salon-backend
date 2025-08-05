package adminTimeSlotTemplate

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminTimeSlotTemplateModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminTimeSlotTemplateService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/time-slot-template"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminTimeSlotTemplateService.GetAllServiceInterface
}

func NewGetAll(service adminTimeSlotTemplateService.GetAllServiceInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Parse query parameters
	var req adminTimeSlotTemplateModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set default limit and offset
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	// parse sort
	var sort []string
	if req.Sort != nil {
		sort = strings.Split(*req.Sort, ",")
	}

	parsedRequest := adminTimeSlotTemplateModel.GetAllParsedRequest{
		Name:   req.Name,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	// Call service
	response, err := h.service.GetAll(c.Request.Context(), parsedRequest)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
