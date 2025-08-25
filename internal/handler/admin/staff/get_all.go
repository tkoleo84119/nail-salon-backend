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

	// trim username, email
	if req.Username != nil {
		*req.Username = strings.TrimSpace(*req.Username)
	}
	if req.Email != nil {
		*req.Email = strings.TrimSpace(*req.Email)
	}

	// set default value
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedRequest := adminStaffModel.GetAllParsedRequest{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: req.IsActive,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	}

	// Service layer call
	response, err := h.service.GetAll(c.Request.Context(), parsedRequest)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
