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

type DeleteBulk struct {
	service adminStoreAccessService.DeleteBulkInterface
}

func NewDeleteBulk(service adminStoreAccessService.DeleteBulkInterface) *DeleteBulk {
	return &DeleteBulk{
		service: service,
	}
}

func (h *DeleteBulk) DeleteBulk(c *gin.Context) {
	// Get target staff ID from path parameter
	staffId := c.Param("staffId")
	if staffId == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"staffId": "staffId為必填項目",
		})
		return
	}
	parsedStaffId, err := utils.ParseID(staffId)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"staffId": "staffId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminStoreAccessModel.DeleteBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	parsedStoreIDs := make([]int64, len(req.StoreIDs))
	for i, storeID := range req.StoreIDs {
		parsedStoreID, err := utils.ParseID(storeID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"storeIds": "storeIds 類型轉換失敗",
			})
			return
		}
		parsedStoreIDs[i] = parsedStoreID
	}

	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Convert store IDs to int64 for permission check
	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.DeleteBulk(c.Request.Context(), parsedStaffId, parsedStoreIDs, staffContext.UserID, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
