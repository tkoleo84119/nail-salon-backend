package adminStaff

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminStaffService.GetAllInterface
}

func NewGetAll(service adminStaffService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Bind query parameters
	var req adminStaffModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// set default value
	limit := 20
	offset := 0
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}
	if req.Offset != nil && *req.Offset >= 0 {
		offset = *req.Offset
	}

	sort := []string{}
	if req.Sort != nil && *req.Sort != "" {
		sort = strings.Split(*req.Sort, ",")
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), adminStaffModel.GetAllParsedRequest{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	})
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
