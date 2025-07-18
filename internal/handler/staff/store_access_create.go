package staff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	staff "github.com/tkoleo84119/nail-salon-backend/internal/model/staff"
	staffService "github.com/tkoleo84119/nail-salon-backend/internal/service/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreAccessHandler struct {
	service staffService.CreateStoreAccessServiceInterface
}

func NewCreateStoreAccessHandler(service staffService.CreateStoreAccessServiceInterface) *CreateStoreAccessHandler {
	return &CreateStoreAccessHandler{service: service}
}

func (h *CreateStoreAccessHandler) CreateStoreAccess(c *gin.Context) {
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Get target staff ID from path parameter
	targetID := c.Param("id")
	if targetID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValInputValidationFailed, map[string]string{
			"id": "id為必填項目",
		})
		return
	}

	// Parse and validate request
	var req staff.CreateStoreAccessRequest
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

	// Extract store IDs from staff context
	var creatorStoreIDs []int64
	for _, store := range staffContext.StoreList {
		storeID, err := utils.ParseID(store.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
			return
		}
		creatorStoreIDs = append(creatorStoreIDs, storeID)
	}

	// Call service
	response, isNewlyCreated, err := h.service.CreateStoreAccess(c.Request.Context(), targetID, req, creatorID, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response with appropriate status code
	if isNewlyCreated {
		c.JSON(http.StatusCreated, common.SuccessResponse(response))
	} else {
		c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}
