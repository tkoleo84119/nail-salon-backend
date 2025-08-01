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

type GetStaffListHandler struct {
	service adminStaffService.GetStaffListServiceInterface
}

func NewGetStaffListHandler(service adminStaffService.GetStaffListServiceInterface) *GetStaffListHandler {
	return &GetStaffListHandler{
		service: service,
	}
}

func (h *GetStaffListHandler) GetStaffList(c *gin.Context) {
	// Bind query parameters
	var req adminStaffModel.GetStaffListRequest
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
	response, err := h.service.GetStaffList(c.Request.Context(), adminStaffModel.GetStaffListParsedRequest{
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
