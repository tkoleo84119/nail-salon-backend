package store

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	storeService "github.com/tkoleo84119/nail-salon-backend/internal/service/store"
)

type GetStoreHandler struct {
	service storeService.GetStoreServiceInterface
}

func NewGetStoreHandler(service storeService.GetStoreServiceInterface) *GetStoreHandler {
	return &GetStoreHandler{
		service: service,
	}
}

func (h *GetStoreHandler) GetStore(c *gin.Context) {
	// Path parameter validation
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
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

	// Service layer call
	response, err := h.service.GetStore(c.Request.Context(), storeID)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}