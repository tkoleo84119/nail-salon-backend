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
	// Input JSON validation
	var req adminStoreModel.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	// Business logic validation
	if !req.HasUpdates() {
		errorCodes.RespondWithEmptyFieldError(c)
		return
	}

	// Authentication context validation
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Service layer call
	response, err := h.service.UpdateStore(c.Request.Context(), storeID, req, *staffContext)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
