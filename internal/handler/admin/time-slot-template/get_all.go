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

type GetTimeSlotTemplateListHandler struct {
	service adminTimeSlotTemplateService.GetTimeSlotTemplateListServiceInterface
}

func NewGetTimeSlotTemplateListHandler(service adminTimeSlotTemplateService.GetTimeSlotTemplateListServiceInterface) *GetTimeSlotTemplateListHandler {
	return &GetTimeSlotTemplateListHandler{
		service: service,
	}
}

func (h *GetTimeSlotTemplateListHandler) GetTimeSlotTemplateList(c *gin.Context) {
	// Parse query parameters
	var req adminTimeSlotTemplateModel.GetTimeSlotTemplateListRequest
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

	parsedRequest := adminTimeSlotTemplateModel.GetTimeSlotTemplateListParsedRequest{
		Name:   req.Name,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	// Call service
	response, err := h.service.GetTimeSlotTemplateList(c.Request.Context(), parsedRequest)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
