package store

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoresHandler struct {
	service storeService.GetStoresServiceInterface
}

func NewGetStoresHandler(service storeService.GetStoresServiceInterface) *GetStoresHandler {
	return &GetStoresHandler{
		service: service,
	}
}

func (h *GetStoresHandler) GetStores(c *gin.Context) {
	// Query parameter validation
	var queryParams storeModel.GetStoresQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Authentication context validation (required per spec)
	_, exists := middleware.GetCustomerFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Set defaults if not provided
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}

	// Service layer call
	response, err := h.service.GetStores(c.Request.Context(), queryParams)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
