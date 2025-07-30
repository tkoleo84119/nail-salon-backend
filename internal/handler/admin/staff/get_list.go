package adminStaff

import (
	"net/http"

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

	// Service layer call
	response, err := h.service.GetStaffList(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
