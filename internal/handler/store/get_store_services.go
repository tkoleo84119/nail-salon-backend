package store

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	storeModel "github.com/tkoleo84119/nail-salon-backend/internal/model/store"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreServicesHandler struct {
	service storeService.GetStoreServicesServiceInterface
}

func NewGetStoreServicesHandler(service storeService.GetStoreServicesServiceInterface) *GetStoreServicesHandler {
	return &GetStoreServicesHandler{
		service: service,
	}
}

func (h *GetStoreServicesHandler) GetStoreServices(c *gin.Context) {
	// Query parameter validation
	var queryParams storeModel.GetStoreServicesQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Path parameter validation
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"storeId": "storeId為必填項目",
		})
		return
	}

	// Set defaults if not provided
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}

	// Service layer call (no authentication required per spec)
	response, err := h.service.GetStoreServices(c.Request.Context(), storeID, queryParams)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
