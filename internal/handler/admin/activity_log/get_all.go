package adminActivityLog

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminActivityLogModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/activity_log"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminActivityLogService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/activity_log"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminActivityLogService.GetAllInterface
}

func NewGetAll(service adminActivityLogService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Parse query parameters
	var req adminActivityLogModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Set default values
	limit := 20
	if req.Limit != nil {
		limit = *req.Limit
	}

	parsedReq := adminActivityLogModel.GetAllParsedRequest{
		Limit: limit,
	}

	// Call service
	response, err := h.service.GetAll(c.Request.Context(), parsedReq)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
