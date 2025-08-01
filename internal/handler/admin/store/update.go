package adminStore

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type UpdateStoreHandler struct {
	service adminStoreService.UpdateStoreServiceInterface
}

func NewUpdateStoreHandler(service adminStoreService.UpdateStoreServiceInterface) *UpdateStoreHandler {
	return &UpdateStoreHandler{
		service: service,
	}
}

func (h *UpdateStoreHandler) UpdateStore(c *gin.Context) {
	// Get store ID from path parameter & parse to int64
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId為必填項目",
		})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"storeId": "storeId 類型轉換失敗"})
		return
	}

	// Input JSON validation
	var req adminStoreModel.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Business logic validation
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// Authentication context validation
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIds := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIds[i], err = utils.ParseID(store.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"storeId": "storeId 類型轉換失敗"})
			return
		}
	}

	// Service layer call
	response, err := h.service.UpdateStore(c.Request.Context(), parsedStoreID, req, staffContext.Role, storeIds)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
