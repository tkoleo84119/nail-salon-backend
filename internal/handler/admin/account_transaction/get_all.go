package adminAccountTransaction

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAccountTransactionService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type GetAll struct {
	service adminAccountTransactionService.GetAllInterface
}

func NewGetAll(service adminAccountTransactionService.GetAllInterface) *GetAll {
	return &GetAll{
		service: service,
	}
}

func (h *GetAll) GetAll(c *gin.Context) {
	// Get store ID from path parameter
	storeID := c.Param("storeId")
	if storeID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"storeId": "storeId 為必填項目",
		})
		return
	}
	parsedStoreID, err := utils.ParseID(storeID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"storeId": "storeId 類型轉換失敗",
		})
		return
	}

	// Get account ID from path parameter
	accountID := c.Param("accountId")
	if accountID == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"accountId": "accountId 為必填項目",
		})
		return
	}
	parsedAccountID, err := utils.ParseID(accountID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"accountId": "accountId 類型轉換失敗",
		})
		return
	}

	// Parse and validate request
	var req adminAccountTransactionModel.GetAllRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Set default values
	limit, offset := utils.SetDefaultValuesOfPagination(req.Limit, req.Offset, 20, 0)
	sort := utils.TransformSort(req.Sort)

	parsedReq := adminAccountTransactionModel.GetAllParsedRequest{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	storeIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		storeIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.GetAll(c.Request.Context(), parsedStoreID, parsedAccountID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
