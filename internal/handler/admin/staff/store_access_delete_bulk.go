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

type DeleteStoreAccessBulkHandler struct {
	service adminStaffService.DeleteStoreAccessBulkServiceInterface
}

func NewDeleteStoreAccessBulkHandler(service adminStaffService.DeleteStoreAccessBulkServiceInterface) *DeleteStoreAccessBulkHandler {
	return &DeleteStoreAccessBulkHandler{
		service: service,
	}
}

func (h *DeleteStoreAccessBulkHandler) DeleteStoreAccessBulk(c *gin.Context) {
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
	var req adminStaffModel.DeleteStoreAccessBulkRequest
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

	// Convert UserID to int64
	creatorID, err := utils.ParseID(staffContext.UserID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"staffId": "staffId 類型轉換失敗",
		})
		return
	}

	// Convert store IDs to int64 for permission check
	var creatorStoreIDs []int64
	for _, storeStr := range staffContext.StoreList {
		storeID, err := utils.ParseID(storeStr.ID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"storeIds": "storeIds 類型轉換失敗",
			})
			return
		}
		creatorStoreIDs = append(creatorStoreIDs, storeID)
	}

	// Call service
	response, err := h.service.DeleteStoreAccessBulk(c.Request.Context(), parsedStaffId, parsedStoreIDs, creatorID, staffContext.Role, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
