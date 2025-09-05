package adminExpense

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminExpenseModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminExpenseService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/expense"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Create struct {
	service adminExpenseService.CreateInterface
}

func NewCreate(service adminExpenseService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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
			"storeId": "storeId 轉換類型失敗",
		})
		return
	}

	var req adminExpenseModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	req.Category = strings.TrimSpace(req.Category)
	if req.Note != nil {
		*req.Note = strings.TrimSpace(*req.Note)
	}

	parsedSupplierID, err := utils.ParseID(req.SupplierID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"supplierId": "supplierId 轉換類型失敗",
		})
		return
	}

	expenseDate, err := utils.DateStringToTime(req.ExpenseDate)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
			"expenseDate": "expenseDate 日期格式錯誤，應為 YYYY-MM-DD",
		})
		return
	}

	var parsedPayerID *int64
	if req.PayerID != nil && *req.PayerID != "" {
		payerID, err := utils.ParseID(*req.PayerID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"payerId": "payerId 轉換類型失敗",
			})
			return
		}
		parsedPayerID = &payerID
	}

	amount := int64(0)
	if req.Amount != nil {
		amount = *req.Amount
	}

	parsedReq := adminExpenseModel.CreateParsedRequest{
		SupplierID:  parsedSupplierID,
		Category:    req.Category,
		Amount:      amount,
		ExpenseDate: expenseDate,
		Note:        req.Note,
		PayerID:     parsedPayerID,
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

	// Call service layer
	response, err := h.service.Create(c.Request.Context(), parsedStoreID, parsedReq, creatorStoreIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
