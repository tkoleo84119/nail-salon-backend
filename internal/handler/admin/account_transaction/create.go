package adminAccountTransaction

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminAccountTransactionModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminAccountTransactionService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/account_transaction"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminAccountTransactionService.CreateInterface
}

func NewCreate(service adminAccountTransactionService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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
	var req adminAccountTransactionModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// trim name, note
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	parsedTransactionDate, err := utils.DateStringToTime(req.TransactionDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"transactionDate": "transactionDate 日期格式錯誤，應為 YYYY-MM-DD",
		})
		return
	}

	parsedReq := adminAccountTransactionModel.CreateParsedRequest{
		TransactionDate: parsedTransactionDate,
		Type:            req.Type,
		Amount:          req.Amount,
		Note:            req.Note,
	}

	// Get staff context from middleware
	staffContext, exists := middleware.GetStaffFromContext(c)
	if !exists {
		errorCodes.AbortWithError(c, errorCodes.AuthContextMissing, nil)
		return
	}

	creatorStoreIDs := make([]int64, len(staffContext.StoreList))
	for i, store := range staffContext.StoreList {
		creatorStoreIDs[i] = store.ID
	}

	// Call service
	response, err := h.service.Create(c.Request.Context(), parsedStoreID, parsedAccountID, parsedReq, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
