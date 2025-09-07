package adminExpenseItem

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	errorCodes "github.com/tkoleo84119/nail-salon-backend/internal/errors"
	"github.com/tkoleo84119/nail-salon-backend/internal/middleware"
	adminExpenseItemModel "github.com/tkoleo84119/nail-salon-backend/internal/model/admin/expense_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/model/common"
	adminExpenseItemService "github.com/tkoleo84119/nail-salon-backend/internal/service/admin/expense_item"
	"github.com/tkoleo84119/nail-salon-backend/internal/utils"
)

type Update struct {
	service adminExpenseItemService.UpdateInterface
}

func NewUpdate(service adminExpenseItemService.UpdateInterface) *Update {
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

	// Parse expenseItemId parameter
	expenseItemIDStr := c.Param("expenseItemId")
	if expenseItemIDStr == "" {
		errorCodes.AbortWithError(c, errorCodes.ValPathParamMissing, map[string]string{
			"expenseItemId": "expenseItemId 為必填項目",
		})
		return
	}
	expenseItemID, err := utils.ParseID(expenseItemIDStr)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"expenseItemId": "expenseItemId 類型轉換失敗",
		})
		return
	}

	// Bind and validate request
	var req adminExpenseItemModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Check if at least one field is provided for update
	if !req.HasUpdates() {
		errorCodes.AbortWithError(c, errorCodes.ValAllFieldsEmpty, nil)
		return
	}

	// Parse and validate optional fields
	var parsedProductID *int64
	if req.ProductID != nil {
		*req.ProductID = strings.TrimSpace(*req.ProductID)
		if *req.ProductID == "" {
			errorCodes.AbortWithError(c, errorCodes.ValFieldNoBlank, map[string]string{
				"productId": "productId 不能為空字串",
			})
			return
		}
		parsed, err := utils.ParseID(*req.ProductID)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
				"productId": "productId 類型轉換失敗",
			})
			return
		}
		parsedProductID = &parsed
	}

	var parsedExpirationDate *time.Time
	if req.ExpirationDate != nil {
		date, err := utils.DateStringToTime(*req.ExpirationDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"expirationDate": "expirationDate 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedExpirationDate = &date
	}

	var parsedArrivalDate *time.Time
	if req.ArrivalDate != nil {
		date, err := utils.DateStringToTime(*req.ArrivalDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"arrivalDate": "arrivalDate 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedArrivalDate = &date
	}

	// Trim string fields
	if req.StorageLocation != nil {
		trimmed := strings.TrimSpace(*req.StorageLocation)
		if trimmed == "" {
			req.StorageLocation = nil
		} else {
			req.StorageLocation = &trimmed
		}
	}

	if req.Note != nil {
		trimmed := strings.TrimSpace(*req.Note)
		if trimmed == "" {
			req.Note = nil
		} else {
			req.Note = &trimmed
		}
	}

	parsedReq := adminExpenseItemModel.UpdateParsedRequest{
		ProductID:       parsedProductID,
		Quantity:        req.Quantity,
		Price:           req.Price,
		ExpirationDate:  parsedExpirationDate,
		IsArrived:       req.IsArrived,
		ArrivalDate:     parsedArrivalDate,
		StorageLocation: req.StorageLocation,
		Note:            req.Note,
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
	response, err := h.service.Update(c.Request.Context(), storeID, expenseID, expenseItemID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, common.SuccessResponse(response))
}
