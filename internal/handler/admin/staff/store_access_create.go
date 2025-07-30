package adminStaff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStaffModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStaffService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/staff"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type CreateStoreAccessHandler struct {
	service adminStaffService.CreateStoreAccessServiceInterface
}

func NewCreateStoreAccessHandler(service adminStaffService.CreateStoreAccessServiceInterface) *CreateStoreAccessHandler {
	return &CreateStoreAccessHandler{service: service}
}

func (h *CreateStoreAccessHandler) CreateStoreAccess(c *gin.Context) {
	// Parse and validate request
	var req adminStaffModel.CreateStoreAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Get target staff ID from path parameter
	targetID := c.Param("staffId")
	if targetID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"staffId": "staffId為必填項目",
		})
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
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
