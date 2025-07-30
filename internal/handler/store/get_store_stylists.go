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

type GetStoreStylistsHandler struct {
	service storeService.GetStoreStylistsServiceInterface
}

func NewGetStoreStylistsHandler(service storeService.GetStoreStylistsServiceInterface) *GetStoreStylistsHandler {
	return &GetStoreStylistsHandler{
		service: service,
	}
}

func (h *GetStoreStylistsHandler) GetStoreStylists(c *gin.Context) {
	// Query parameter validation
	var queryParams storeModel.GetStoreStylistsQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Path parameter validation
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId為必填項目",
		})
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
	response, err := h.service.GetStoreStylists(c.Request.Context(), storeID, queryParams)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
