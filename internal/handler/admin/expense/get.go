package adminExpense

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminExpenseService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Get struct {
	service adminExpenseService.GetInterface
}

func NewGet(service adminExpenseService.GetInterface) *Get {
	return &Get{
		service: service,
	}
}

func (h *Get) Get(c *gin.Context) {
	// Parse storeId parameter
	storeIDStr := c.Param("storeId")
	if storeIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	storeID, err := utils.ParseID(storeIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	// Parse expenseId parameter
	expenseIDStr := c.Param("expenseId")
	if expenseIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"expenseId": "expenseId 為必填項目",
		})
		return
	}
	expenseID, err := utils.ParseID(expenseIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"expenseId": "expenseId 類型轉換失敗",
		})
		return
	}

	// Get staff context from middleware
	staff, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	// Extract store IDs for permission check
	storeIDs := make([]int64, len(staff.StoreList))
	for i, store := range staff.StoreList {
		storeIDs[i] = store.ID
	}

	// Call service layer
	response, err := h.service.Get(c.Request.Context(), storeID, expenseID, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
