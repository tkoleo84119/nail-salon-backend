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

type CreateStoreHandler struct {
	service adminStoreService.CreateStoreServiceInterface
}

func NewCreateStoreHandler(service adminStoreService.CreateStoreServiceInterface) *CreateStoreHandler {
	return &CreateStoreHandler{
		service: service,
	}
}

func (h *CreateStoreHandler) CreateStore(c *gin.Context) {
	// Parse and validate request
	var req adminStoreModel.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	staffId, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		// if staffId can't be parsed, it means the staffId is invalid
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"staffId": "staffId 類型轉換失敗"})
		return
	}

	// Call service
	response, err := h.service.CreateStore(c.Request.Context(), req, staffId, staffContext.Role)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
