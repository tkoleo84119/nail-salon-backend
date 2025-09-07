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

type Create struct {
	service adminExpenseItemService.CreateInterface
}

func NewCreate(service adminExpenseItemService.CreateInterface) *Create {
	return &Create{
		service: service,
	}
}

func (h *Create) Create(c *gin.Context) {
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
	var req adminExpenseItemModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := utils.ExtractValidationErrors(err)
		errorCodes.RespondWithValidationErrors(c, validationErrors)
		return
	}

	// Parse and validate required fields
	parsedProductID, err := utils.ParseID(req.ProductID)
	if err != nil {
		errorCodes.AbortWithError(c, errorCodes.ValTypeConversionFailed, map[string]string{
			"productId": "productId 類型轉換失敗",
		})
		return
	}

	// Parse and validate optional fields
	parsedExpirationDate := time.Time{}
	if req.ExpirationDate != nil {
		date, err := utils.DateStringToTime(*req.ExpirationDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"expirationDate": "expirationDate 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedExpirationDate = date
	}

	parsedArrivalDate := time.Time{}
	if req.ArrivalDate != nil {
		date, err := utils.DateStringToTime(*req.ArrivalDate)
		if err != nil {
			errorCodes.AbortWithError(c, errorCodes.ValFieldDateFormat, map[string]string{
				"arrivalDate": "arrivalDate 日期格式錯誤，應為 YYYY-MM-DD",
			})
			return
		}
		parsedArrivalDate = date
	}

	// Trim string fields
	if req.StorageLocation != nil {
		trimmed := strings.TrimSpace(*req.StorageLocation)
		req.StorageLocation = &trimmed
	}
	if req.Note != nil {
		trimmed := strings.TrimSpace(*req.Note)
		req.Note = &trimmed
	}

	quantity := int64(0)
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	price := int64(0)
	if req.Price != nil {
		price = *req.Price
	}

	isArrived := false
	if req.IsArrived != nil {
		isArrived = *req.IsArrived
	}

	parsedReq := adminExpenseItemModel.CreateParsedRequest{
		ProductID:       parsedProductID,
		Quantity:        quantity,
		Price:           price,
		ExpirationDate:  &parsedExpirationDate,
		IsArrived:       isArrived,
		ArrivalDate:     &parsedArrivalDate,
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
	response, err := h.service.Create(c.Request.Context(), storeID, expenseID, parsedReq, storeIDs)
	if err != nil {
		errorCodes.RespondWithServiceError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, common.SuccessResponse(response))
}
