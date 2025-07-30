package adminStore

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	adminStoreModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStoreService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetStoreListHandler struct {
	service adminStoreService.GetStoreListServiceInterface
}

func NewGetStoreListHandler(service adminStoreService.GetStoreListServiceInterface) *GetStoreListHandler {
	return &GetStoreListHandler{
		service: service,
	}
}

func (h *GetStoreListHandler) GetStoreList(c *gin.Context) {
	// Parse and validate query parameters
	var req adminStoreModel.GetStoreListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Service layer call
	response, err := h.service.GetStoreList(c.Request.Context(), req)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
