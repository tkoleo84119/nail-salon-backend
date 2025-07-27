package adminStore

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
)

type GetStoreHandler struct {
	service adminStoreService.GetStoreServiceInterface
}

func NewGetStoreHandler(service adminStoreService.GetStoreServiceInterface) *GetStoreHandler {
	return &GetStoreHandler{
		service: service,
	}
}

func (h *GetStoreHandler) GetStore(c *gin.Context) {
	// Get store ID from path parameter
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{"storeId": "storeId 為必填項目"})
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