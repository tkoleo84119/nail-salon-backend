package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/service"
	serviceService "github.com/tkoleo84119/nail-salon-backend/internal/service/service"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateServiceHandler struct {
	service serviceService.CreateServiceInterface
}

func NewCreateServiceHandler(service serviceService.CreateServiceInterface) *CreateServiceHandler {
	return &CreateServiceHandler{
		service: service,
	}
}

func (h *CreateServiceHandler) CreateService(c *gin.Context) {
	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.RespondWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	response, err := h.service.CreateService(c.Request.Context(), req, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
