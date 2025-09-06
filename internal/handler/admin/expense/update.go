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

type Update struct {
	service adminExpenseService.UpdateInterface
}

func NewUpdate(service adminExpenseService.UpdateInterface) *Update {
	return &Update{
		service: service,
	}
}

func (h *Update) Update(c *gin.Context) {
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

	// Bind and validate request
	var req adminExpenseModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Check if at least one field is provided for update
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, map[string]string{
			"fields": "至少需要提供一個欄位進行更新",
		})
		return
	}

	// Trim string fields
	if req.Category != nil {
		trimmed := strings.TrimSpace(*req.Category)
		req.Category = &trimmed
	}

	if req.Note != nil {
		trimmed := strings.TrimSpace(*req.Note)
		req.Note = &trimmed
	}

	// Parse and validate dates
	parsedReq := adminExpenseModel.UpdateParsedRequest{
		SupplierID:   nil,
		Category:     req.Category,
		Amount:       req.Amount,
		OtherFee:     req.OtherFee,
		ExpenseDate:  nil,
		Note:         req.Note,
		PayerID:      nil,
		IsReimbursed: req.IsReimbursed,
		ReimbursedAt: nil,
	}

	// Parse supplierID
	if req.SupplierID != nil && *req.SupplierID != "" {
		supplierID, err := utils.ParseID(*req.SupplierID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"supplierId": "supplierId 類型轉換失敗",
			})
			return
		}
		parsedReq.SupplierID = &supplierID
	}

	// Parse payerID
	if req.PayerID != nil && *req.PayerID != "" {
		payerID, err := utils.ParseID(*req.PayerID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"payerId": "payerId 類型轉換失敗",
			})
			return
		}
		parsedReq.PayerID = &payerID
	} else if *req.PayerID == "" {
		trueValue := true
		parsedReq.PayerIDIsNone = &trueValue
		parsedReq.PayerID = nil
	}

	// Parse expenseDate
	if req.ExpenseDate != nil && *req.ExpenseDate != "" {
		expenseDate, err := utils.DateStringToTime(*req.ExpenseDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"expenseDate": "expenseDate 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedReq.ExpenseDate = &expenseDate
	}

	// Parse reimbursedAt
	if req.ReimbursedAt != nil && *req.ReimbursedAt != "" {
		reimbursedAt, err := utils.DateStringToTime(*req.ReimbursedAt)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"reimbursedAt": "reimbursedAt 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedReq.ReimbursedAt = &reimbursedAt
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
	response, err := h.service.Update(c.Request.Context(), storeID, expenseID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
