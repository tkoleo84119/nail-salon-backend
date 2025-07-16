package storeAccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/store-access"
	storeAccessService "github.com/tkoleo84119/nail-salon-backend/internal/service/store-access"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type DeleteStoreAccessHandler struct {
	service storeAccessService.DeleteStoreAccessServiceInterface
}

func NewDeleteStoreAccessHandler(service storeAccessService.DeleteStoreAccessServiceInterface) *DeleteStoreAccessHandler {
	return &DeleteStoreAccessHandler{service: service}
}

func (h *DeleteStoreAccessHandler) DeleteStoreAccess(c *gin.Context) {
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Get target staff ID from path parameter
	targetID := c.Param("id")
	if targetID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"id": "員工ID為必填項目",
		})
		return
	}

	// Parse and validate request
	var req storeAccess.DeleteStoreAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, validationErrors)
		return
	}

	// Convert UserID to int64
	creatorID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Convert store IDs to int64 for permission check
	var creatorStoreIDs []int64
	for _, storeStr := range staffContext.StoreList {
		storeID, err := utils.ParseID(storeStr.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
			return
		}
		creatorStoreIDs = append(creatorStoreIDs, storeID)
	}

	// Call service
	response, err := h.service.DeleteStoreAccess(c.Request.Context(), targetID, req, creatorID, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
