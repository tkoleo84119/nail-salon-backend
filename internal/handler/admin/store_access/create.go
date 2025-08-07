package adminStoreAccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminStoreAccessModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/store_access"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminStoreAccessService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/store_access"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminStoreAccessService.CreateInterface
}

func NewCreate(service adminStoreAccessService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
	// Get target staff ID from path parameter
	staffID := c.Param("staffId")
	if staffID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"staffId": "staffId為必填項目",
		})
		return
	}
	parsedStaffID, err := utils.ParseID(staffID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{"staffUserId": "staffUserId 類型轉換失敗"})
		return
	}

	// Parse and validate request
	var req adminStoreAccessModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	parsedStoreID, err := utils.ParseID(req.StoreID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Extract store IDs from staff context
	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	// Call service
	response, isNewlyCreated, err := h.service.Create(c.Request.Context(), parsedStaffID, parsedStoreID, staffContext.UserID, staffContext.Role, creatorStoreIDs)
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
